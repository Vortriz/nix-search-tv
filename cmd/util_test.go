package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"
	"github.com/alecthomas/assert/v2"
)

type state struct {
	CacheDir  string
	ConfigDir string
	Stdout    *bytes.Buffer
}

func setup(t *testing.T) state {
	cacheDir, err := os.MkdirTemp("", "nix-search-tv-cache")
	assert.NoError(t, err)
	err = os.Setenv("XDG_CACHE_HOME", cacheDir)
	assert.NoError(t, err)

	configDir, err := os.MkdirTemp("", "nix-search-tv-config")
	assert.NoError(t, err)
	err = os.Setenv("XDG_CONFIG_HOME", configDir)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	Stdout = buf

	indices.SetFetchers(map[string]indexer.Fetcher{})

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(cacheDir))
		assert.NoError(t, os.RemoveAll(configDir))

		// Set these to "." so that is anything leaks,
		// we'll see in the test directory
		err = os.Setenv("XDG_CACHE_HOME", "./tmp-tests/cache")
		assert.NoError(t, err)
		err = os.Setenv("XDG_CONFIG_HOME", "./tmp-tests/config")
		assert.NoError(t, err)

		Stdout = nil

		indices.SetFetchers(map[string]indexer.Fetcher{})
	})

	return state{
		CacheDir:  cacheDir,
		ConfigDir: configDir,
		Stdout:    buf,
	}
}

func setMetadata(t *testing.T, state state, md indexer.IndexMetadata) {
	mdbytes, err := json.Marshal(md)
	assert.NoError(t, err)

	path := filepath.Join(
		state.CacheDir,
		"nix-search-tv",
		indices.Nixpkgs,
		"metadata.json",
	)
	err = os.WriteFile(path, mdbytes, 0666)
	assert.NoError(t, err)
}

func getCache(t *testing.T, state state) []string {
	path := filepath.Join(
		state.CacheDir,
		"nix-search-tv",
		indices.Nixpkgs,
		"cache.txt",
	)
	cacheb, err := os.ReadFile(path)
	assert.NoError(t, err)

	cache := strings.TrimSpace(string(cacheb))
	return strings.Split(cache, "\n")
}

func setNixpkgs(pkgs ...string) {
	indices.SetFetchers(map[string]indexer.Fetcher{
		indices.Nixpkgs: &PkgsFetcher{
			pkgs: pkgs,
		},
	})
}

type FailFetcher struct{}

func (f *FailFetcher) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	return "", errors.New("failed to get latest release")
}

func (f *FailFetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	return nil, errors.New("failed to download the release")
}

type PkgsFetcher struct {
	pkgs []string
}

func (f *PkgsFetcher) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	return "latest", nil
}

func (f *PkgsFetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	pkgs := indexer.Indexable{Packages: map[string]json.RawMessage{}}
	for _, pkg := range f.pkgs {
		pkgs.Packages[pkg] = []byte("{}")
	}

	data, err := json.Marshal(pkgs)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewBuffer(data)), nil
}
