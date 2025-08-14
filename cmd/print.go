package cmd

import (
	"bufio"
	"context"
	"fmt"
	"maps"
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
	indexNames := conf.Indexes

	customIndexes, err := RegisterCustomIndexes(conf)
	if err != nil {
		return fmt.Errorf("register fetchers: %w", err)
	}

	indexNames = append(indexNames, customIndexes...)

	// If a custom index is passed via flag e.g.
	// `--indexes nvf`, it will appear twice in
	// `indexNames`, which will lead to problems
	// later
	indexesSet := map[string]struct{}{}
	for _, indexName := range indexNames {
		indexesSet[indexName] = struct{}{}
	}
	indexNames = slices.Collect(maps.Keys(indexesSet))

	indexes, err := GetIndexes(conf.CacheDir, indexNames)
	if err != nil {
		return fmt.Errorf("get indexes: %w", err)
	}

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

	withPrefix := len(indexes) > 1

	for _, index := range indexes {
		canPrint := !slices.ContainsFunc(needIndexing, func(need indexer.Index) bool {
			return need.Name == index.Name
		})
		if canPrint {
			err = PrintIndexKeys(conf, index.Name, withPrefix)
			if err != nil {
				return fmt.Errorf("%s: %w", index, err)
			}
		}
	}

	results := indexer.RunIndexing(ctx, conf.CacheDir, needIndexing)
	for result := range results {
		if result.Err != nil {
			msg := addIndexPrefix(
				result.Index,
				fmt.Sprintf("indexing failed: %s\n", result.Err),
			)
			Stdout.Write([]byte(msg))
			continue
		}

		err := PrintIndexKeys(conf, result.Index, withPrefix)
		if err != nil {
			return fmt.Errorf("%s: %w", result.Index, err)
		}
	}

	return nil
}

func PrintIndexKeys(conf config.Config, index string, withPrefix bool) error {
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
	if withPrefix {
		prefix = []byte(index + "/ ")
	}

	slices.Sort(allkeys)

	for _, k := range allkeys {
		Stdout.Write(append([]byte(prefix), []byte(k)...))
		Stdout.Write([]byte{'\n'})
	}

	return nil
}
