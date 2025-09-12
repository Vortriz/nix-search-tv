package cmd

import (
	"github.com/3timeslazy/nix-search-tv/indexes/indices"
	"github.com/urfave/cli/v3"
)

var Homepage = &cli.Command{
	Name:      "homepage",
	UsageText: "nix-search-tv homepage [package_name]",
	Usage:     "Print the link to the package homepage",
	Action:    NewPreviewAction(indices.HomepagePreview),
	Flags:     BaseFlags(),
}
