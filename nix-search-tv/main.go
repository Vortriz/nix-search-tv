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
		IndexCmd(ctx, searcher)
		return
	}
	if cmd == "print" {
		Print(os.Stdout)
		return
	}
	if cmd == "preview" {
		if len(os.Args) != 3 {
			fmt.Println("A package path is required. Ex. nixpkgs.nix-search")
			return
		}

		pkgPath := os.Args[2]
		pkg, err := SearchPkg(ctx, searcher, pkgPath)
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
		return fmt.Errorf("failed to read %s: %w", nixpkgsFile, err)
	}

	_, err = wr.Write(asBytes)
	return err
}

func defaultIndexPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot get user cache dir: %w", err)
	}

	return filepath.Join(cacheDir, "nix-search-tv"), nil
}

func generatePkgsList(ctx context.Context, out io.Writer, searcher *blugesearcher.PackagesSearcher) error {
	pkgsSeq, err := searcher.SearchPackages(ctx, "", search.Opts{})
	if err != nil {
		return err
	}
	for pkg := range pkgsSeq {
		fmt.Fprintln(out, pkg.Path)
	}
	return nil
}

func SearchPkg(ctx context.Context, searcher *blugesearcher.PackagesSearcher, pkgPath string) (search.SearchedPackage, error) {
	query, _ := strings.CutPrefix(pkgPath, "nixpkgs.")
	pkgs, err := searcher.SearchPackages(ctx, query, search.Opts{
		Exact: true,
	})
	if err != nil {
		return search.SearchedPackage{}, fmt.Errorf("nix-search error: %w", err)
	}

	for pkg := range pkgs {
		if pkg.Path == pkgPath {
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

	fmt.Fprint(out, styler.Bold("license"), "\n")
	if pkg.Unfree {
		fmt.Fprint(out, "unfree $$$")
	} else {
		fmt.Fprint(out, "free as in freedom \\o/")
	}
	fmt.Fprint(out, "\n\n")

	if pkg.MainProgram != "" {
		fmt.Fprint(out, styler.Bold("main program"), "\n")
		fmt.Fprint(out, style.PrintCodeBlock("$ "+pkg.MainProgram))
		fmt.Fprint(out, "\n\n")
	}
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
