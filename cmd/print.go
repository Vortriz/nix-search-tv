package cmd

import (
	"bufio"
	"context"
	"fmt"
	"slices"
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
	indexes := conf.Indexes

	registered, err := RegisterRenderDocs(conf)
	if err != nil {
		return fmt.Errorf("register render-docs: %w", err)
	}
	indexes = append(indexes, registered...)

	needIndexing, err := indexer.NeedIndexing(
		conf.CacheDir,
		time.Duration(conf.UpdateInterval),
		indexes,
	)
	if err != nil {
		return fmt.Errorf("check if indexing needed: %w", err)
	}

	if len(needIndexing) > 0 {
		if conf.EnableWaitingMessage {
			PrintWaiting(Stdout)
		}
	}

	for _, index := range indexes {
		if !slices.Contains(needIndexing, index) {
			err = PrintIndexKeys(index, conf)
			if err != nil {
				return fmt.Errorf("%s: %w", index, err)
			}
		}
	}
	err = Index(ctx, conf, needIndexing)
	if err != nil {
		return fmt.Errorf("index packages: %w", err)
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
		allkeys = append(allkeys, scanner.Text())
	}

	prefix := []byte{}
	if len(conf.Indexes) > 1 {
		prefix = []byte(index + "/ ")
	}

	slices.Sort(allkeys)

	for _, k := range allkeys {
		Stdout.Write(append([]byte(prefix), []byte(k)...))
		Stdout.Write([]byte{'\n'})
	}

	return nil
}
