package indexer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"time"

	"github.com/3timeslazy/nix-search-tv/config"
)

type Fetcher interface {
	GetLatestRelease(context.Context, IndexMetadata) (string, error)
	DownloadRelease(context.Context, string) (io.ReadCloser, error)
}

type Index struct {
	Name    string
	Fetcher Fetcher
}

type IndexMetadata struct {
	LastIndexedAt time.Time `json:"last_indexed_at"`
	CurrRelease   string    `json:"curr_release"`
}

func RunIndexing(
	ctx context.Context,
	conf config.Config,
	indexes []Index,
) []error {
	errs := make([]error, len(indexes))

	wg := sync.WaitGroup{}

	for i, index := range indexes {
		indexDir := filepath.Join(conf.CacheDir, index.Name)
		md, err := GetMetadata(indexDir)
		if err != nil {
			errs[i] = fmt.Errorf("get metadata: %w", err)
			continue
		}
		if time.Since(md.LastIndexedAt) < time.Duration(conf.UpdateInterval) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			indexDir := filepath.Join(conf.CacheDir, index.Name)
			md, err := GetMetadata(indexDir)
			if err != nil {
				errs[i] = fmt.Errorf("get metadata: %w", err)
				return
			}

			latest, err := index.Fetcher.GetLatestRelease(ctx, md)
			if err != nil {
				errs[i] = fmt.Errorf("get latest release: %w", err)
				return
			}
			if latest == md.CurrRelease {
				_ = SetMetadata(indexDir, IndexMetadata{
					LastIndexedAt: time.Now(),
					CurrRelease:   latest,
				})
				return
			}

			pkgs, err := index.Fetcher.DownloadRelease(ctx, latest)
			if err != nil {
				errs[i] = fmt.Errorf("download latest release: %w", err)
				return
			}
			cache, err := CacheWriter(indexDir)
			if err != nil {
				errs[i] = err
				return
			}

			badgerDir := filepath.Join(indexDir, "badger")
			indexer, err := NewBadger(badgerDir)
			if err != nil {
				errs[i] = fmt.Errorf("open indexer: %w", err)
				return
			}

			err = indexer.Index(pkgs, cache)
			if err != nil {
				errs[i] = fmt.Errorf("index packages: %w", err)
				return
			}

			_ = SetMetadata(indexDir, IndexMetadata{
				LastIndexedAt: time.Now(),
				CurrRelease:   latest,
			})
		}()
	}
	wg.Wait()

	return errs
}

const Nixpkgs = "nixpkgs"

var ErrUnknownIndex = errors.New("unknown index")
