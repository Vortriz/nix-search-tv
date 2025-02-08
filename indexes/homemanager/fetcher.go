package homemanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexer"
)

type Fetcher struct{}

func (f *Fetcher) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	cmd := exec.Command(
		"nix", "build",
		"github:nix-community/home-manager/master#docs-json",
		"--no-write-lock-file", "--no-link", "--print-out-paths",
	)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("build options: %w", err)
	}

	path := strings.TrimSpace(string(out))
	return filepath.Join(path, "/share/doc/home-manager/options.json"), nil
}

func (f *Fetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	file, err := os.OpenFile(release, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("open options file: %w", err)
	}

	return newOptionsWrapper(file), nil
}

type optionsWrapper struct {
	opts io.Closer
	mrd  io.Reader
}

// newOptionsWrapper translates home-manager's format into the
// format of the indexer
//
// home-manager:
//
//	{
//		"pkg1": { ... },
//		"pkg2": { ... }
//	}
//
// indexer:
//
//	{
//		"packages": {
//		  "pkg1": { ... }
//		}
//	}
func newOptionsWrapper(rd io.ReadCloser) *optionsWrapper {
	mrd := io.MultiReader(
		bytes.NewBufferString(`{"packages":`),
		rd,
		bytes.NewBufferString(`}`),
	)
	return &optionsWrapper{
		opts: rd,
		mrd:  mrd,
	}
}

func (ow *optionsWrapper) Read(p []byte) (n int, err error) {
	return ow.mrd.Read(p)
}

func (ow *optionsWrapper) Close() error {
	return ow.opts.Close()
}
