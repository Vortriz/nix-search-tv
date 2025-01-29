package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/3timeslazy/nix-search-tv/nix-search-tv/style"

	"libdb.so/nix-search/search"
	"libdb.so/nix-search/search/searchers/blugesearcher"
)

var ErrPkgNotFound = errors.New("package not found")

const availableCommands = `Options are "index", "print", "preview"`

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("A command is required. %s\n", availableCommands)
		return
	}

	ctx := context.Background()
	searcher, err := blugesearcher.Open("")
	if err != nil {
		fmt.Printf("nix-search error: %v.\nHave you run nix-search?\n", err)
		return
	}
	defer searcher.Close()

	cmd := os.Args[1]
	if cmd == "index" {
		err := IndexCmd(ctx, searcher)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	if cmd == "print" {
		err := Print(os.Stdout)
		if err != nil {
			fmt.Printf("failed to print package list: %v\n", err)
		}
		return
	}
	if cmd == "preview" {
		if len(os.Args) != 3 {
			fmt.Println("A package path is required. Ex. nixpkgs.nix-search")
			return
		}

		fullPkgName := os.Args[2]
		pkg, err := SearchPkg(ctx, searcher, fullPkgName)
		if err != nil {
			if errors.Is(err, ErrPkgNotFound) {
				fmt.Println("package not found")
				return
			}

			fmt.Println(err)
			return
		}

		PreviewPkg(os.Stdout, pkg)
		return
	}

	fmt.Printf("Unknown command %q. %s\n", os.Args[1], availableCommands)
	return
}

const indexFile = "nixpkgs.txt"

func IndexCmd(ctx context.Context, searcher *blugesearcher.PackagesSearcher) error {
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

func Print(wr io.Writer) error {
	indexDir, err := defaultIndexPath()
	if err != nil {
		return err
	}

	nixpkgsFile := filepath.Join(indexDir, indexFile)
	asBytes, err := os.ReadFile(nixpkgsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("index file not found. Run `nix-search-tv index` and try again")
		}
		return fmt.Errorf("failed to read %s: %w", nixpkgsFile, err)
	}

	_, err = wr.Write(asBytes)
	return err
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

func generatePkgsList(ctx context.Context, out io.Writer, searcher *blugesearcher.PackagesSearcher) error {
	pkgsSeq, err := searcher.SearchPackages(ctx, "", search.Opts{})
	if err != nil {
		return err
	}
	for pkg := range pkgsSeq {
		fullPkgName := CutChannel(pkg.Path)
		fmt.Fprintln(out, fullPkgName)
	}
	return nil
}

func SearchPkg(ctx context.Context, searcher *blugesearcher.PackagesSearcher, fullPkgName string) (search.SearchedPackage, error) {
	query := CutSubpkg(fullPkgName)
	pkgs, err := searcher.SearchPackages(ctx, query, search.Opts{
		Exact: true,
	})
	if err != nil {
		return search.SearchedPackage{}, fmt.Errorf("nix-search error: %w", err)
	}

	for pkg := range pkgs {
		if CutChannel(pkg.Path) == fullPkgName {
			return pkg, nil
		}
	}

	return search.SearchedPackage{}, ErrPkgNotFound
}

func PreviewPkg(out io.Writer, pkg search.SearchedPackage) {
	styler := style.StyledText

	// title
	fmt.Fprint(out, styler.Bold(pkg.Name))
	fmt.Fprint(out, " ", styler.Dim("("+pkg.Version+")"), "\n")

	fmt.Fprint(out, style.Wrap(pkg.Description, ""))
	// two new lines instead of one here and after to make `tv` render it as a single new line
	fmt.Fprint(out, "\n\n")

	if pkg.LongDescription != "" && pkg.Description != pkg.LongDescription {
		fmt.Fprint(out, style.StyleLongDescription(styler, pkg.LongDescription), "\n")
	}

	if len(pkg.Homepages) > 0 {
		subtitle := "homepage"
		if len(pkg.Homepages) > 1 {
			subtitle += "s"
		}
		fmt.Fprint(out, styler.Bold(subtitle), "\n")
		fmt.Fprint(out, strings.Join(pkg.Homepages, "\n"))
		fmt.Fprint(out, "\n\n")
	}

	licenses := strings.Join(pkg.Licenses, "\n")
	licenseType := "free"
	if pkg.Unfree {
		licenseType = "unfree"
	}
	fmt.Fprint(out, styler.Bold("license"))
	fmt.Fprint(out, styler.Dim(" ("+licenseType+")"), "\n")
	fmt.Fprint(out, licenses)
	fmt.Fprint(out, "\n\n")

	if pkg.MainProgram != "" {
		fmt.Fprint(out, styler.Bold("main program"), "\n")
		fmt.Fprint(out, style.PrintCodeBlock("$ "+pkg.MainProgram))
		fmt.Fprint(out, "\n\n")
	}
}
