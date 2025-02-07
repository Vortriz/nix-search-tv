package cmd

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/nixpkgs"

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
	flagIndexer := []string{Nixpkgs}

	indexes, mds, err := indexer.NeedIndexing(conf, flagIndexer)
	if err != nil {
		return fmt.Errorf("check if indexing needed: %w", err)
	}
	if len(mds) > 0 {
		if conf.EnableWaitingMessage.Bool {
			PrintWaiting(os.Stdout)
		}

		err = Index(ctx, conf, indexes)
		if err != nil {
			return err
		}
	}

	for _, index := range flagIndexer {
		keys, err := indexer.OpenKeysReader(conf.CacheDir, index)
		if err != nil {
			return fmt.Errorf("failed to read %s keys: %w", index, err)
		}
		defer keys.Close()

		if _, err := io.Copy(os.Stdout, keys); err != nil {
			return err
		}
	}

	return nil
}

const (
	Nixpkgs = "nixpkgs"
)

var ErrUnknownIndex = errors.New("unknown index")

func Index(ctx context.Context, conf config.Config, indexNames []string) error {
	indexes := []indexer.Index{}
	for _, indexName := range indexNames {
		switch indexName {
		case Nixpkgs:
			indexes = append(indexes, indexer.Index{
				Name:    indexName,
				Fetcher: &nixpkgs.Fetcher{},
			})

		default:
			return ErrUnknownIndex
		}
	}

	errs := indexer.RunIndexing(ctx, conf, indexes)
	success := false
	for _, err := range errs {
		if err == nil {
			success = true
			continue
		}
	}
	if !success {
		return fmt.Errorf("all indexes failed: %w", errs[0])
	}
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
