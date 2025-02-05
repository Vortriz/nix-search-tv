package cmd

import (
	"context"
	"os"

	"github.com/3timeslazy/nix-search-tv/tv"
	"github.com/urfave/cli/v3"
)

var Print = &cli.Command{
	Name:      "print",
	UsageText: "nix-search-tv print",
	Usage:     "Print the list of all index Nix packages\nSupposed to be called by Television",
	Action:    PrintAction,
}

func PrintAction(_ context.Context, _ *cli.Command) error {
	return tv.PrintPackages(os.Stdout)
}
