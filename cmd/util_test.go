package cmd

import (
	"bytes"
	"context"
	"encoding/json"
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

	setNixpkgs("nix-search-tv")

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

		setNixpkgs("nix-search-tv")
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

	path := filepath.Join(state.CacheDir, indices.Nixpkgs, "metadata.json")
	err = os.WriteFile(path, mdbytes, 0666)
	assert.NoError(t, err)
}

func getCache(t *testing.T, state state) []string {
	cacheb, err := os.ReadFile(filepath.Join(state.CacheDir, indices.Nixpkgs, "cache.txt"))
	assert.NoError(t, err)

	cache := strings.TrimSpace(string(cacheb))
	return strings.Split(cache, "\n")
}

func setNixpkgs(pkgs ...string) {
	Fetchers = map[string]indexer.Fetcher{
		indices.Nixpkgs: &NixpkgsFetcher{
			pkgs: pkgs,
		},
	}
}

type NixpkgsFetcher struct {
	pkgs []string
}

func (f *NixpkgsFetcher) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	return indices.Nixpkgs, nil
}

func (f *NixpkgsFetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	pkgs := struct {
		Packages map[string]json.RawMessage `json:"packages"`
	}{
		Packages: map[string]json.RawMessage{},
	}
	for _, pkg := range f.pkgs {
		pkgs.Packages[pkg] = []byte("{}")
	}

	data, err := json.Marshal(pkgs)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewBuffer(data)), nil
}
