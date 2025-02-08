package cmd

import (
	"cmp"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexes"
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
				if !indexes.Indexes[ind] {
					avail := slices.Collect(maps.Keys(indexes.Indexes))
					return fmt.Errorf("unknown index %q. Available options are: %v", ind, avail)
				}
			}
			return nil
		},
	}
)

func GetConfig(cmd *cli.Command) (config.Config, error) {
	var err error

	conf := config.Default()
	path := cmd.String(ConfigFlag.Name)
	if path == "" {
		path, err = config.ConfigDir()
		if err != nil {
			return conf, fmt.Errorf("get default config path: %w", err)
		}
	}
	loaded, err := config.LoadPath(path)
	if err != nil {
		return conf, fmt.Errorf("load config: %w", err)
	}

	conf.UpdateInterval = cmp.Or(loaded.UpdateInterval, conf.UpdateInterval)
	conf.CacheDir = cmp.Or(loaded.CacheDir, conf.CacheDir)

	if loaded.Indexes.Valid {
		conf.Indexes = loaded.Indexes
	}
	if loaded.EnableWaitingMessage.Valid {
		conf.EnableWaitingMessage = loaded.EnableWaitingMessage
	}

	if err := os.MkdirAll(conf.CacheDir, 0755); err != nil {
		return conf, fmt.Errorf("cannot create cache directory: %w", err)
	}

	return conf, nil
}
