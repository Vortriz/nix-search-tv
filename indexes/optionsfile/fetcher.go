package optionsfile

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/readutil"
)

type Fetcher struct {
	path string
}

var _ indexer.OptionFileFetcher = (*Fetcher)(nil)

func NewFetcher(path string) *Fetcher {
	return &Fetcher{
		path: path,
	}
}

func (fetcher *Fetcher) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	// It is expected to be a nix store path, which already contains a hash in its name.
	// So, it the options file haven't changed since last build, the store path name should remain the same.
	// Otherwise, it will be different and trigger indexing
	return fetcher.path, nil
}

func (fetcher *Fetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	file, err := os.Open(fetcher.path)
	if err != nil {
		return nil, fmt.Errorf("open packages file: %w", err)
	}

	return readutil.PackagesWrapper(file), nil
}

func (fetcher *Fetcher) Path() string {
	return fetcher.path
}
