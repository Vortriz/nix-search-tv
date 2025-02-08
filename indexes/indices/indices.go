package indices

import (
	"fmt"
	"os"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/homemanager"
	"github.com/3timeslazy/nix-search-tv/indexes/nixpkgs"
)

const (
	Nixpkgs     = "nixpkgs"
	HomeManager = "home-manager"
)

var Indexes = map[string]bool{
	Nixpkgs:     true,
	HomeManager: true,
}

var Fetchers = map[string]indexer.Fetcher{
	Nixpkgs:     &nixpkgs.Fetcher{},
	HomeManager: &homemanager.Fetcher{},
}

var Previews = map[string]func(config.Config, string) error{
	Nixpkgs: func(conf config.Config, pkgName string) error {
		pkg, err := indexer.LoadKey[nixpkgs.Package](conf, Nixpkgs, pkgName)
		if err != nil {
			return fmt.Errorf("load package content: %w", err)
		}

		nixpkgs.Preview(os.Stdout, pkg)
		return nil
	},
	HomeManager: func(conf config.Config, pkgName string) error {
		pkg, err := indexer.LoadKey[homemanager.Package](conf, HomeManager, pkgName)
		if err != nil {
			return fmt.Errorf("load package content: %w", err)
		}

		homemanager.Preview(os.Stdout, pkg)
		return nil
	},
}
