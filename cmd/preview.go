package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexes/indices"

	"github.com/urfave/cli/v3"
)

var Preview = &cli.Command{
	Name:      "preview",
	UsageText: "nix-search-tv preview [package_name]",
	Usage:     "Print preview for the package",
	Action:    PreviewAction,
	Flags:     baseFlags,
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

	inds := cmd.StringSlice(IndexesFlag.Name)
	if len(inds) == 0 {
		inds = conf.Indexes
	}

	if len(inds) == 1 {
		preview := indices.Previews[inds[0]]
		return preview(conf, fullPkgName)
	}

	ind, pkgName, ok := strings.Cut(fullPkgName, ":")
	if !ok {
		return errors.New("multiple indexes requested, but the package has no index prefix")
	}

	preview := indices.Previews[ind]
	return preview(conf, pkgName)
}
