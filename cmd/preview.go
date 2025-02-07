package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/nixpkgs"
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
	if fullPkgName == waitingMessage {
		PreviewWaiting(os.Stdout, conf)
		return nil
	}

	pkg, err := indexer.LoadKey[nixpkgs.Package](conf, Nixpkgs, fullPkgName)
	if err != nil {
		return fmt.Errorf("load package: %w", err)
	}

	nixpkgs.Preview(os.Stdout, pkg)
	return nil
}
