package tv

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"libdb.so/nix-search/search"
	"libdb.so/nix-search/search/searchers/blugesearcher"
)

const indexFile = "nixpkgs.txt"

func GeneratePackagesFile(ctx context.Context) error {
	searcher, err := blugesearcher.Open("")
	if err != nil {
		return fmt.Errorf("open nix-search searcher: %w", err)
	}
	defer searcher.Close()

	indexDir, err := defaultIndexPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(indexDir, 0755); err != nil {
		return fmt.Errorf("cannot create index directory: %w", err)
	}

	nixpkgsFile, err := os.Create(filepath.Join(indexDir, indexFile))
	if err != nil {
		return fmt.Errorf("cannot create nixpkgs.txt: %w", err)
	}
	err = generatePkgsList(ctx, nixpkgsFile, searcher)
	if err != nil {
		return fmt.Errorf("nix-search failed: %w", err)
	}

	return nil
}

func generatePkgsList(ctx context.Context, out io.Writer, searcher *blugesearcher.PackagesSearcher) error {
	pkgsSeq, err := searcher.SearchPackages(ctx, "", search.Opts{})
	if err != nil {
		return err
	}
	for pkg := range pkgsSeq {
		fullPkgName := cutChannel(pkg.Path)
		fmt.Fprintln(out, fullPkgName)
	}
	return nil
}

func defaultIndexPath() (string, error) {
	var err error
	// Check xdg first because `os.UserCacheDir`
	// ignores XDG_CACHE_HOME on darwin
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		cacheDir, err = os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("cannot get user cache dir: %w", err)
		}
	}

	return filepath.Join(cacheDir, "nix-search-tv"), nil
}
