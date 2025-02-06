package cmd

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/metafiles"
	"github.com/3timeslazy/nix-search-tv/nixpkgs"
	"github.com/3timeslazy/nix-search-tv/nixpkgs/indexer"

	"github.com/urfave/cli/v3"
)

var Print = &cli.Command{
	Name:      "print",
	UsageText: "nix-search-tv print",
	Usage:     "Print the list of all index Nix packages",
	Action:    PrintAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "path to the configuration file",
		},
	},
}

func PrintAction(ctx context.Context, cmd *cli.Command) error {
	conf, err := GetConfig(cmd)
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}
	md, err := metafiles.GetMetadata[nixpkgs.Metadata](conf.CacheDir)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}

	if time.Since(md.LastIndexedAt) > time.Duration(conf.UpdateInterval) {
		if conf.EnableWaitingMessage.Bool {
			PrintWaiting(os.Stdout)
		}

		err = Index(ctx, conf, md.CurrRelease)
		if err != nil {
			return fmt.Errorf("failed to index: %w", err)
		}
	}

	cache, err := metafiles.CacheReader(conf.CacheDir)
	if err != nil {
		return fmt.Errorf("failed read cache: %w", err)
	}

	_, _ = io.Copy(os.Stdout, cache)
	return nil
}

func Index(ctx context.Context, conf config.Config, currRelease string) error {
	indexerDir := path.Join(conf.CacheDir, "badger")
	indexer, err := indexer.NewBadger(indexerDir)
	if err != nil {
		return fmt.Errorf("create indexer: %w", err)
	}
	defer indexer.Close()

	release, err := nixpkgs.FindLatestRelease(ctx, currRelease)
	if err != nil {
		return fmt.Errorf("find latest release: %w", err)
	}
	if currRelease == release {
		// Don't check the error, because displayed
		// and outdated results are better than nothing
		_ = metafiles.SetMetadata(conf.CacheDir, nixpkgs.Metadata{
			LastIndexedAt: time.Now(),
			CurrRelease:   release,
		})
		return nil
	}

	pkgs, err := nixpkgs.DownloadRelease(ctx, release)
	if err != nil {
		return fmt.Errorf("download release: %w", err)
	}
	defer pkgs.Close()

	cache, err := metafiles.CacheWriter(conf.CacheDir)
	if err != nil {
		return err
	}
	defer cache.Close()

	err = indexer.Index(pkgs, cache)
	if err != nil {
		return fmt.Errorf("index packages: %w", err)
	}

	_ = metafiles.SetMetadata(conf.CacheDir, nixpkgs.Metadata{
		LastIndexedAt: time.Now(),
		CurrRelease:   release,
	})
	return nil
}

func GetConfig(cmd *cli.Command) (config.Config, error) {
	var err error

	conf := config.Default()
	path := cmd.String("config")
	if path == "" {
		path, err = config.ConfigDir()
		if err != nil {
			return conf, fmt.Errorf("get default config path: %w", err)
		}
	}
	loaded, err := config.LoadPath(path)
	if err != nil {
		return conf, fmt.Errorf("load config: %w", err)
	}

	conf.UpdateInterval = cmp.Or(loaded.UpdateInterval, conf.UpdateInterval)
	conf.CacheDir = cmp.Or(loaded.CacheDir, conf.CacheDir)

	if loaded.EnableWaitingMessage.Valid {
		conf.EnableWaitingMessage = loaded.EnableWaitingMessage
	}

	if err := os.MkdirAll(conf.CacheDir, 0755); err != nil {
		return conf, fmt.Errorf("cannot create cache directory: %w", err)
	}

	return conf, nil
}
