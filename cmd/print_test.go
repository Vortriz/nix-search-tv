package cmd

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"
	"github.com/alecthomas/assert/v2"
	"github.com/urfave/cli/v3"
)

func TestPrint_Errors(t *testing.T) {
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

func TestPrint_Config(t *testing.T) {
	ctx := context.Background()

	Fetchers = map[string]indexer.Fetcher{
		indices.Nixpkgs: &NixpkgsFetcher{},
	}

	runPrint := func(t *testing.T, ctx context.Context, args []string) {
		cmd := cli.Command{
			Flags:  BaseFlags(),
			Action: PrintAction,
		}
		err := cmd.Run(ctx, append([]string{"print"}, args...))
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

		runPrint(t, ctx, nil)

		data := state.Stdout.String()
		assert.NotContains(t, data, waitingMessage)
		assert.Contains(t, data, "nix-search-tv")
	})

	t.Run("use defaults when default file not exist", func(t *testing.T) {
		state := setup(t)

		runPrint(t, ctx, nil)

		// enable_waiting_message is true by default, so
		// if we see it in the output - the defaults were used
		data := state.Stdout.String()
		assert.Contains(t, data, waitingMessage)
		assert.Contains(t, data, "nix-search-tv")
	})
}

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
	})

	return state{
		CacheDir:  cacheDir,
		ConfigDir: configDir,
		Stdout:    buf,
	}
}

type NixpkgsFetcher struct{}

func (f *NixpkgsFetcher) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	return "nixpkgs", nil
}

func (f *NixpkgsFetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	pkgs := bytes.NewBufferString(`{
	  "packages": {
	    "nix-search-tv": {}
	  }
	}`)

	return io.NopCloser(pkgs), nil
}

// func (f *NixpkgsFetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
// 	pkgs := bytes.NewBufferString(`
// {
//   "packages": {
//     "nix-search-tv": {
//       "meta": {
//         "available": true,
//         "broken": false,
//         "changelog": "https://github.com/3timeslazy/nix-search-tv/releases/tag/v1.0.0",
//         "description": "Nixpkgs channel for television",
//         "homepage": "https://github.com/3timeslazy/nix-search-tv",
//         "insecure": false,
//         "license": {
//           "deprecated": false,
//           "free": true,
//           "fullName": "GNU General Public License v3.0 only",
//           "redistributable": true,
//           "shortName": "gpl3Only",
//           "spdxId": "GPL-3.0-only",
//           "url": "https://spdx.org/licenses/GPL-3.0-only.html"
//         },
//         "mainProgram": "nix-search-tv",
//         "maintainers": [
//           {
//             "email": "gaetan@glepage.com",
//             "github": "GaetanLepage",
//             "githubId": 33058747,
//             "name": "Gaetan Lepage"
//           }
//         ],
//         "name": "nix-search-tv-1.0.0",
//         "outputsToInstall": [
//           "out"
//         ],
//         "platforms": [
//           "x86_64-darwin",
//           "i686-darwin",
//           "aarch64-darwin",
//           "armv7a-darwin",
//           "aarch64-linux",
//           "armv5tel-linux",
//           "armv6l-linux",
//           "armv7a-linux",
//           "armv7l-linux",
//           "i686-linux",
//           "loongarch64-linux",
//           "m68k-linux",
//           "microblaze-linux",
//           "microblazeel-linux",
//           "mips-linux",
//           "mips64-linux",
//           "mips64el-linux",
//           "mipsel-linux",
//           "powerpc64-linux",
//           "powerpc64le-linux",
//           "riscv32-linux",
//           "riscv64-linux",
//           "s390-linux",
//           "s390x-linux",
//           "x86_64-linux",
//           "wasm64-wasi",
//           "wasm32-wasi",
//           "i686-freebsd",
//           "x86_64-freebsd",
//           "aarch64-freebsd"
//         ],
//         "position": "pkgs/by-name/ni/nix-search-tv/package.nix:33",
//         "unfree": false,
//         "unsupported": false
//       },
//       "name": "nix-search-tv-1.0.0",
//       "outputName": "out",
//       "outputs": {
//         "out": null
//       },
//       "pname": "nix-search-tv",
//       "system": "x86_64-linux",
//       "version": "1.0.0"
//     }
//   }
// }
// 	`)

// 	return io.NopCloser(pkgs), nil
// }
