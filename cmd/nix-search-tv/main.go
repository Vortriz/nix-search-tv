package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/3timeslazy/nix-search-tv/tv"
	"libdb.so/nix-search/search"

	"github.com/hashicorp/go-hclog"
	"github.com/urfave/cli/v3"
)

var opts = search.DefaultIndexPackageOpts

var app = cli.App{
	Name:      "nix-search-tv",
	UsageText: `nix-search-tv [options] [command]`,
	Usage:     "A tool integrating nix-search and television",
	Commands: []*cli.Command{
		IndexCmd,
		{
			Name:      "print",
			UsageText: "nix-search-tv print",
			Usage:     "Print the list of all index Nix packages\nSupposed to be called by Television",
			Action:    PrintCmd,
		},
		{
			Name:      "preview",
			UsageText: "nix-search-tv preview [package_name]",
			Usage:     "Print preview for the package",
			Action:    PreviewCmd,
		},
	},
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		code := 1

		var codeError cli.ExitCoder
		if errors.As(err, &codeError) {
			code = codeError.ExitCode()
		}

		log := hclog.FromContext(ctx)
		log.Error("error", "err", err)

		os.Exit(code)
	}
}

func PrintCmd(c *cli.Context) error {
	return tv.PrintPackages(os.Stdout)
}

func PreviewCmd(c *cli.Context) error {
	fullPkgName := c.Args().First()
	if fullPkgName == "" {
		return errors.New("package name is required")
	}

	return tv.PrintPreview(c.Context, os.Stdout, fullPkgName)
}
