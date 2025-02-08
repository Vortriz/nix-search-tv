package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"

	"github.com/urfave/cli/v3"
)

var (
	baseFlags = []cli.Flag{
		ConfigFlag,
		IndexesFlag,
	}

	ConfigFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "path to the configuration file",
	}

	IndexesFlag = &cli.StringSliceFlag{
		Name:  "indexes",
		Usage: "what packages to index",
		Validator: func(indexNames []string) error {
			for _, ind := range indexNames {
				if !indices.Indexes[ind] {
					avail := slices.Collect(maps.Keys(indices.Indexes))
					return fmt.Errorf("unknown index %q. Available options are: %v", ind, avail)
				}
			}
			return nil
		},
	}
)

func GetConfig(cmd *cli.Command) (config.Config, error) {
	path := cmd.String(ConfigFlag.Name)
	conf, err := config.LoadPath(path)
	if err != nil {
		return config.Config{}, fmt.Errorf("load config: %w", err)
	}

	if err := os.MkdirAll(conf.CacheDir, 0755); err != nil {
		return conf, fmt.Errorf("cannot create cache directory: %w", err)
	}

	return conf, nil
}
