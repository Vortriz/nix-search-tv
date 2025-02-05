package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/3timeslazy/nix-search-tv/tv"
	"github.com/urfave/cli/v3"
)

var Preview = &cli.Command{
	Name:      "preview",
	UsageText: "nix-search-tv preview [package_name]",
	Usage:     "Print preview for the package",
	Action:    PreviewCmd,
}

func PreviewCmd(ctx context.Context, cmd *cli.Command) error {
	fullPkgName := cmd.Args().First()
	if fullPkgName == "" {
		return errors.New("package name is required")
	}

	return tv.PrintPreview(ctx, os.Stdout, fullPkgName)
}
