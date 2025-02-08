package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/3timeslazy/nix-search-tv/indexer"

	"github.com/urfave/cli/v3"
)

var Print = &cli.Command{
	Name:      "print",
	UsageText: "nix-search-tv print",
	Usage:     "Print the list of all index Nix packages",
	Action:    PrintAction,
	Flags:     baseFlags,
}

func PrintAction(ctx context.Context, cmd *cli.Command) error {
	conf, err := GetConfig(cmd)
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}
	indexes := cmd.StringSlice(IndexesFlag.Name)
	if len(indexes) == 0 {
		indexes = conf.Indexes.V
	}

	needIndexing, mds, err := indexer.NeedIndexing(conf, indexes)
	if err != nil {
		return fmt.Errorf("check if indexing needed: %w", err)
	}
	if len(mds) > 0 {
		if conf.EnableWaitingMessage.Bool {
			PrintWaiting(os.Stdout)
		}

		err = Index(ctx, conf, needIndexing)
		if err != nil {
			return err
		}
	}

	needPrefix := len(indexes) > 1
	for _, index := range indexes {
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

func PrintKeys(prefix string, pkgs io.Reader) error {
	if prefix == "" {
		_, err := io.Copy(os.Stdout, pkgs)
		return err
	}

	prefixb := []byte(prefix + ":")
	scanner := bufio.NewScanner(pkgs)
	for scanner.Scan() {
		os.Stdout.Write(append(prefixb, append(scanner.Bytes(), '\n')...))
	}

	return nil
}
