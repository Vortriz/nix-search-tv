package cmd

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"

	"github.com/urfave/cli/v3"
)

func BaseFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  ConfigFlag,
			Usage: "path to the configuration file",
			Validator: func(path string) error {
				if path == "" {
					return errors.New("config path cannot be empty")
				}
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name:  IndexesFlag,
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
		},
		&cli.StringFlag{
			Name:   CacheDirFlag,
			Hidden: true,
			Usage:  "Path to the indexes cache directory",
		},
	}
}

const (
	ConfigFlag   = "config"
	IndexesFlag  = "indexes"
	CacheDirFlag = "cache-dir"
)

var Stdout io.ReadWriter = os.Stdout

func GetConfig(cmd *cli.Command) (config.Config, error) {
	var conf config.Config
	var err error

	if cmd.IsSet(ConfigFlag) {
		conf, err = config.LoadPath(cmd.String(ConfigFlag))
	} else {
		conf, err = config.LoadDefault()
	}
	if err != nil {
		return config.Config{}, fmt.Errorf("load config: %w", err)
	}

	if cmd.IsSet(IndexesFlag) {
		conf.Indexes = cmd.StringSlice(IndexesFlag)
	}
	if cmd.IsSet(CacheDirFlag) {
		conf.CacheDir = cmd.String(CacheDirFlag)
	}

	if err := os.MkdirAll(conf.CacheDir, 0755); err != nil {
		return conf, fmt.Errorf("cannot create cache directory: %w", err)
	}

	return conf, nil
}
