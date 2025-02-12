package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"

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

	needIndexing, mds, err := indexer.NeedIndexing(conf, conf.Indexes)
	if err != nil {
		return fmt.Errorf("check if indexing needed: %w", err)
	}
	if len(mds) > 0 {
		if conf.EnableWaitingMessage {
			PrintWaiting(Stdout)
		}

		err = Index(ctx, conf, needIndexing)
		if err != nil {
			return err
		}
	}

	needPrefix := len(conf.Indexes) > 1
	for _, index := range conf.Indexes {
		keys, err := indexer.OpenKeysReader(conf.CacheDir, index)
		if err != nil {
			return fmt.Errorf("failed to read %s keys: %w", index, err)
		}
		defer keys.Close()

		prefix := ""
		if needPrefix {
			prefix = index
		}
		if err = PrintKeys(prefix, keys); err != nil {
			return err
		}
	}

	return nil
}

func PrintKeys(index string, pkgs io.Reader) error {
	if index == "" {
		_, err := io.Copy(Stdout, pkgs)
		return err
	}

	scanner := bufio.NewScanner(pkgs)
	for scanner.Scan() {
		out := addIndexPrefix(index, scanner.Text()+"\n")
		Stdout.Write([]byte(out))
	}

	return nil
}
