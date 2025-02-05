package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/3timeslazy/nix-search-tv/nixpkgs"
	"github.com/3timeslazy/nix-search-tv/nixpkgs/indexer"
	"github.com/urfave/cli/v3"
)

var Preview = &cli.Command{
	Name:      "preview",
	UsageText: "nix-search-tv preview [package_name]",
	Usage:     "Print preview for the package",
	Action:    PreviewAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "path to the configuration file",
		},
	},
}

func PreviewAction(ctx context.Context, cmd *cli.Command) error {
	fullPkgName := cmd.Args().First()
	if fullPkgName == "" {
		return errors.New("package name is required")
	}

	conf, err := GetConfig(cmd)
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	indexerDir := path.Join(conf.CacheDir, "badger")
	indexer, err := indexer.NewBadger(indexerDir)
	if err != nil {
		return fmt.Errorf("open indexer: %w", err)
	}
	defer indexer.Close()

	pkg, err := indexer.Load(fullPkgName)
	if err != nil {
		return fmt.Errorf("load package: %w", err)
	}

	nixpkgs.Preview(os.Stdout, pkg)
	return nil
}
