package cmd

import (
	"bufio"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexer"

	"github.com/urfave/cli/v3"
)

var Print = &cli.Command{
	Name:      "print",
	UsageText: "nix-search-tv print",
	Usage:     "Print the list of all index Nix packages",
	Action:    PrintAction,
	Flags:     BaseFlags(),
}

func PrintAction(ctx context.Context, cmd *cli.Command) error {
	conf, err := GetConfig(cmd)
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	needIndexing, err := indexer.NeedIndexing(
		conf.CacheDir,
		time.Duration(conf.UpdateInterval),
		conf.Indexes,
	)
	if err != nil {
		return fmt.Errorf("check if indexing needed: %w", err)
	}

	if len(needIndexing) > 0 {
		if conf.EnableWaitingMessage {
			PrintWaiting(Stdout)
		}
	}

	for _, index := range conf.Indexes {
		if !slices.Contains(needIndexing, index) {
			err = PrintIndexKeys(index, conf)
			if err != nil {
				return fmt.Errorf("%s: %w", index, err)
			}
		}
	}
	err = Index(ctx, conf, needIndexing)
	if err != nil {
		return err
	}

	return nil
}

func PrintIndexKeys(index string, conf config.Config) error {
	keys, err := indexer.OpenKeysReader(conf.CacheDir, index)
	if err != nil {
		return fmt.Errorf("read keys file: %w", err)
	}
	defer keys.Close()

	allkeys := []string{}
	scanner := bufio.NewScanner(keys)
	for scanner.Scan() {
		if len(conf.Indexes) == 1 {
			allkeys = append(allkeys, scanner.Text())
		} else {
			allkeys = append(allkeys, addIndexPrefix(index, scanner.Text()))
		}
	}

	slices.Sort(allkeys)
	Stdout.Write([]byte(strings.Join(allkeys, "\n")))

	return nil
}
