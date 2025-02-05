package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/3timeslazy/nix-search-tv/cmd"
	"libdb.so/nix-search/search"

	"github.com/hashicorp/go-hclog"
	"github.com/urfave/cli/v3"
)

var opts = search.DefaultIndexPackageOpts

var root = &cli.Command{
	Name:      "nix-search-tv",
	UsageText: `nix-search-tv [options] [command]`,
	Usage:     "Nix-related television channel",
	Commands: []*cli.Command{
		cmd.Print,
		cmd.Preview,
	},
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := root.Run(ctx, os.Args); err != nil {
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
