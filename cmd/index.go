package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"
)

var ErrUnknownIndex = errors.New("unknown index")

func Index(ctx context.Context, conf config.Config, indexNames []string) error {
	indexes := []indexer.Index{}
	for _, indexName := range indexNames {
		fetcher, ok := indices.GetFetcher(indexName)
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownIndex, indexName)
		}
		indexes = append(indexes, indexer.Index{
			Name:    indexName,
			Fetcher: fetcher,
		})
	}

	results := indexer.RunIndexing(ctx, conf.CacheDir, indexes)
	for result := range results {
		if result.Err != nil {
			msg := addIndexPrefix(
				result.Index,
				fmt.Sprintf("indexing failed: %s\n", result.Err),
			)
			Stdout.Write([]byte(msg))
			continue
		}

		err := PrintIndexKeys(result.Index, conf)
		if err != nil {
			return err
		}
	}

	return nil
}

// The two functions below connect the print and preview
// commands. Their logic is simple, so the only reason
// these functions exist is to keep prefix logic in one place
//
// Also, pay attention to the fact the `addIndexPrefix` puts a
// space between index and package names, while `cutIndexPrefix`
// cut without the space. That's done this way because when there is
// a space, tv and fzf consider it as two arguments. Those two arguments
// later passed by tv/fzf into the preview command as "{1}{2}" and that's
// where the space dissappears

func addIndexPrefix(index, pkg string) string {
	return index + "/ " + pkg
}

func cutIndexPrefix(s string) (string, string, bool) {
	index, pkg, ok := strings.Cut(s, "/")
	return strings.TrimSpace(index), strings.TrimSpace(pkg), ok
}
