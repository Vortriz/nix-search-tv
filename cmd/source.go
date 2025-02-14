package cmd

import (
	"github.com/3timeslazy/nix-search-tv/indexes/indices"
	"github.com/urfave/cli/v3"
)

var Source = &cli.Command{
	Name:      "source",
	UsageText: "nix-search-tv source [package_name]",
	Usage:     "Print the the link to the source code",
	Action:    NewPreviewAction(indices.SourcePreview),
	Flags:     BaseFlags(),
}
