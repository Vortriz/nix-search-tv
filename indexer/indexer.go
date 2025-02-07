package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	Name     string
	Fetcher  Fetcher
	Metadata IndexMetadata
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
	wg.Add(len(indexes))

	for i, index := range indexes {
		go func() {
			defer wg.Done()

			indexDir := filepath.Join(conf.CacheDir, index.Name)
			latest, err := index.Fetcher.GetLatestRelease(ctx, index.Metadata)
			if err != nil {
				errs[i] = fmt.Errorf("get latest release: %w", err)
				return
			}
			if latest == index.Metadata.CurrRelease {
				_ = setIndexMetadata(indexDir, IndexMetadata{
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
			defer pkgs.Close()

			cache, err := CacheWriter(indexDir)
			if err != nil {
				errs[i] = err
				return
			}
			defer cache.Close()

			badgerDir := filepath.Join(indexDir, "badger")
			indexer, err := NewBadger(badgerDir)
			if err != nil {
				errs[i] = fmt.Errorf("open indexer: %w", err)
				return
			}
			defer indexer.Close()

			err = indexer.Index(pkgs, cache)
			if err != nil {
				errs[i] = fmt.Errorf("index packages: %w", err)
				return
			}

			_ = setIndexMetadata(indexDir, IndexMetadata{
				LastIndexedAt: time.Now(),
				CurrRelease:   latest,
			})
		}()
	}

	wg.Wait()

	return errs
}

func NeedIndexing(
	conf config.Config,
	indexes []string,
) ([]string, []IndexMetadata, error) {
	mds := []IndexMetadata{}
	needIndex := []string{}

	for _, index := range indexes {
		indexDir := filepath.Join(conf.CacheDir, index)
		md, err := getIndexMetadata(indexDir)
		if err != nil {
			return nil, nil, fmt.Errorf("get metadata: %w", err)
		}
		if time.Since(md.LastIndexedAt) > time.Duration(conf.UpdateInterval) {
			mds = append(mds, md)
			needIndex = append(needIndex, index)
		}
	}

	return needIndex, mds, nil
}

func OpenKeysReader(cacheDir, index string) (io.ReadCloser, error) {
	indexDir := filepath.Join(cacheDir, index)
	path, err := initFile(indexDir, cacheFile, nil)
	if err != nil {
		return nil, fmt.Errorf("init cache file: %w", err)
	}

	return os.OpenFile(path, os.O_RDONLY, 0666)
}

func LoadKey[T any](conf config.Config, index, key string) (T, error) {
	var pkg T

	badgerDir := filepath.Join(conf.CacheDir, index, "badger")
	indexer, err := NewBadger(badgerDir)
	if err != nil {
		return pkg, fmt.Errorf("open indexer: %w", err)
	}
	defer indexer.Close()

	data, err := indexer.Load(key)
	if err != nil {
		return pkg, fmt.Errorf("load key: %w", err)
	}

	err = json.Unmarshal(data, &pkg)
	if err != nil {
		return pkg, fmt.Errorf("decode package: %w", err)
	}

	return pkg, nil
}
