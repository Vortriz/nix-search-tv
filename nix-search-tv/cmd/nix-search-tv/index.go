package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"libdb.so/nix-search/search"
	"libdb.so/nix-search/search/searchers/blugesearcher"

	"github.com/3timeslazy/nix-search-tv/nix-search-tv/tv"
	"github.com/hashicorp/go-hclog"
	"github.com/urfave/cli/v3"
)

var IndexCmd = &cli.Command{
	Name:   "index",
	Usage:  "update the Nixpkgs index before searching",
	Action: Index,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "index-path",
			Usage:   "path to the index directory, defaults to a directory in $XDG_CACHE_HOME",
			EnvVars: []string{"NIX_SEARCH_INDEX_PATH"},
		},
		&cli.StringFlag{
			Name:        "channel",
			Aliases:     []string{"c"},
			Usage:       "channel path to index",
			Value:       opts.Nixpkgs,
			Destination: &opts.Nixpkgs,
			Action: func(ctx *cli.Context, v string) error {
				if !strings.HasPrefix(v, "<") || !strings.HasSuffix(v, ">") {
					return errors.New(fmt.Sprintf("invalid channel %q", v))
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:  "flake",
			Usage: "flake to index unless channel is provided",
			Action: func(c *cli.Context, v string) error {
				path, err := search.ResolveNixPathFromFlake(c.Context, c.String("flake"))
				if err != nil {
					return fmt.Errorf("failed to resolve flake: %w", err)
				}
				c.Set("flake", path)
				return nil
			},
		},
		&cli.IntFlag{
			Name:        "max-jobs",
			Aliases:     []string{"j"},
			Usage:       "max parallel jobs",
			Value:       opts.Parallelism,
			Destination: &opts.Parallelism,
		},
	},
}

const indexFile = "nixpkgs.txt"

func Index(c *cli.Context) error {
	ctx := c.Context
	log := hclog.FromContext(ctx)
	indexPath := c.String("index-path")

	if !blugesearcher.Exists(indexPath) {
		log.Info("first run or outdated index detected, will index packages")
	}

	if c.IsSet("flake") {
		if c.IsSet("channel") {
			return errors.New("cannot set both --channel and --flake")
		}

		opts.Nixpkgs = c.String("flake")
	}

	log.Info("indexing packages")

	if err := indexNixSearch(ctx, indexPath); err != nil {
		return err
	}
	if err := tv.GeneratePackagesFile(ctx); err != nil {
		return fmt.Errorf("failed to generate package file: %w", err)
	}

	return nil
}

func indexNixSearch(ctx context.Context, indexPath string) error {
	pkgs, err := search.IndexPackages(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to get package index: %w", err)
	}

	if err := blugesearcher.IndexPackages(ctx, indexPath, pkgs); err != nil {
		return fmt.Errorf("failed to store indexed packages: %w", err)
	}

	return nil
}

func defaultIndexPath() (string, error) {
	var err error
	// Check xdg first because `os.UserCacheDir`
	// ignores XDG_CACHE_HOME on darwin
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		cacheDir, err = os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("cannot get user cache dir: %w", err)
		}
	}

	return filepath.Join(cacheDir, "nix-search-tv"), nil
}
