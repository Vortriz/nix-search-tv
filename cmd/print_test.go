package cmd

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"
	"github.com/alecthomas/assert/v2"
	"github.com/urfave/cli/v3"
)

func TestPrintInvalidFlags(t *testing.T) {
	ctx := context.Background()

	runPrint := func(args ...string) error {
		cmd := cli.Command{
			Writer: io.Discard,
			Flags:  BaseFlags(),
			Action: PrintAction,
		}
		return cmd.Run(ctx, append([]string{"print"}, args...))
	}

	t.Run("given config not exist", func(t *testing.T) {
		setup(t)

		err := runPrint("--config", "unknown.json")
		assert.IsError(t, err, fs.ErrNotExist)
		assert.Contains(t, err.Error(), "unknown.json")
	})

	t.Run("unknown index", func(t *testing.T) {
		setup(t)

		err := runPrint("--indexes", "nixpkgs,unknown")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown index")
	})
}

func TestPrintConfig(t *testing.T) {
	ctx := context.Background()

	runPrint := func(t *testing.T, ctx context.Context) {
		cmd := cli.Command{
			Writer: io.Discard,
			Flags:  BaseFlags(),
			Action: PrintAction,
		}
		err := cmd.Run(ctx, []string{"print"})
		assert.NoError(t, err)
	}

	t.Run("read default file if no config given", func(t *testing.T) {
		state := setup(t)

		confdir := filepath.Join(state.ConfigDir, "nix-search-tv")
		assert.NoError(t, os.MkdirAll(confdir, 0755))

		// enable_waiting_message is true by default, so
		// if we don't see it in the output - the config file was used
		conffile := filepath.Join(confdir, "config.json")
		assert.NoError(t, os.WriteFile(conffile, []byte(`{ "enable_waiting_message": false }`), 0666))

		runPrint(t, ctx)

		data := state.Stdout.String()
		assert.NotContains(t, data, waitingMessage)
		assert.Contains(t, data, "nix-search-tv")
	})

	t.Run("use defaults when default file not exist", func(t *testing.T) {
		state := setup(t)

		runPrint(t, ctx)

		// enable_waiting_message is true by default, so
		// if we see it in the output - the defaults were used
		data := state.Stdout.String()
		assert.Contains(t, data, waitingMessage)
		assert.Contains(t, data, "nix-search-tv")
	})
}

func TestPrintIndexing(t *testing.T) {
	ctx := context.Background()

	runPrint := func(t *testing.T, ctx context.Context, cacheDir string) {
		cmd := cli.Command{
			Writer: io.Discard,
			Flags:  BaseFlags(),
			Action: PrintAction,
		}
		err := cmd.Run(ctx, append([]string{"print"}, []string{
			"--indexes", indices.Nixpkgs,
			"--cache-dir", cacheDir,
		}...))
		assert.NoError(t, err)
	}

	t.Run("run -> no index -> indexing", func(t *testing.T) {
		state := setup(t)

		runPrint(t, ctx, state.CacheDir)

		expectedPaths := map[string]bool{
			"nixpkgs":               false,
			"nixpkgs/badger":        false,
			"nixpkgs/cache.txt":     false,
			"nixpkgs/metadata.json": false,
		}
		err := filepath.WalkDir(state.CacheDir, func(path string, d fs.DirEntry, err error) error {
			for expectedPath := range expectedPaths {
				if strings.HasSuffix(path, expectedPath) {
					expectedPaths[expectedPath] = true
				}
			}
			return nil
		})
		assert.NoError(t, err)

		for expectedPath, ok := range expectedPaths {
			assert.True(t, ok, "File not found: %s", expectedPath)
		}

		cache := getCache(t, state)
		assertSortEqual(t, []string{"nix-search-tv"}, cache)
	})

	t.Run("has index -> need indexing -> indexing", func(t *testing.T) {
		state := setup(t)

		// generate the first index
		runPrint(t, ctx, state.CacheDir)

		// reset metadata and run again
		newpkgs := []string{"televison", "fzf"}
		setNixpkgs(newpkgs...)
		setMetadata(t, state, indexer.IndexMetadata{})
		runPrint(t, ctx, state.CacheDir)

		cache := getCache(t, state)
		assertSortEqual(t, newpkgs, cache)
	})

	t.Run("has index -> indexing not needed", func(t *testing.T) {
		state := setup(t)

		// generate the first index
		runPrint(t, ctx, state.CacheDir)

		// set fetchers to nil. If everything is correct, they won't be triggered. If
		// there's a bug, the test will panic
		indices.SetFetchers(nil)
		assert.NotPanics(t, func() {
			runPrint(t, ctx, state.CacheDir)
		})
	})

	t.Run("single index failed", func(t *testing.T) {
		state := setup(t)

		setFailingFetcher()
		runPrint(t, ctx, state.CacheDir)

		expected := fmt.Sprintf("%s/ indexing failed", indices.Nixpkgs)
		assert.Contains(t, state.Stdout.String(), expected)
	})

	t.Run("need update, but not new version", func(t *testing.T) {

	})

	t.Run("run multiple indexes at once", func(t *testing.T) {

	})
}

func assertSortEqual[S ~[]E, E cmp.Ordered](t *testing.T, expected, actual S) {
	slices.Sort(expected)
	slices.Sort(actual)
	assert.Equal(t, expected, actual)
}
