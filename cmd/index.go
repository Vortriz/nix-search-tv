package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes"
)

var ErrUnknownIndex = errors.New("unknown index")

func Index(ctx context.Context, conf config.Config, indexNames []string) error {
	indexesToRun := []indexer.Index{}

	for _, indexName := range indexNames {
		fetcher, ok := indexes.Fetchers[indexName]
		if !ok {
			return ErrUnknownIndex
		}
		indexesToRun = append(indexesToRun, indexer.Index{
			Name:    indexName,
			Fetcher: fetcher,
		})
	}

	errs := indexer.RunIndexing(ctx, conf, indexesToRun)
	success := false
	for _, err := range errs {
		if err == nil {
			success = true
			continue
		}
	}
	if !success {
		return fmt.Errorf("all indexes failed: %w", errs[0])
	}
	return nil
}
