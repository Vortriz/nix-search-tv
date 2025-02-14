package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"

	"github.com/urfave/cli/v3"
)

var Preview = &cli.Command{
	Name:      "preview",
	UsageText: "nix-search-tv preview [package_name]",
	Usage:     "Print preview for the package",
	Action:    NewPreviewAction(indices.Preview),
	Flags:     BaseFlags(),
}

type PreviewFunc func(out io.Writer, index string, pkg json.RawMessage) error

func NewPreviewAction(preview PreviewFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		fullPkgName := strings.Join(cmd.Args().Slice(), " ")
		if fullPkgName == "" {
			return errors.New("package name is required")
		}

		conf, err := GetConfig(cmd)
		if err != nil {
			return fmt.Errorf("get config: %w", err)
		}
		if fullPkgName == waitingMessage {
			PreviewWaiting(Stdout, conf)
			return nil
		}

		var index, pkgName string

		if len(conf.Indexes) == 1 {
			index = conf.Indexes[0]
			pkgName = fullPkgName
		} else {
			var ok bool
			index, pkgName, ok = cutIndexPrefix(fullPkgName)
			if !ok {
				return errors.New("multiple indexes requested, but the package has no index prefix")
			}
		}

		pkg, err := indexer.LoadKey(conf, index, pkgName)
		if err != nil {
			return fmt.Errorf("load package content: %w", err)
		}

		return preview(Stdout, index, pkg)
	}
}
