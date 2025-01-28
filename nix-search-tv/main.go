package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
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
	searcher := Must(blugesearcher.Open(""))
	defer searcher.Close()

	cmd := os.Args[1]
	if cmd == "index" {
		GeneratePkgsList(ctx, os.Stdout, searcher)
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

func Print(wr io.Writer) error {
	// TODO: save to xdg cache
	absPath := "/home/vladimir/kode/src/github.com/3timeslazy/nix-search-tv/nix-search-tv/pkgs.txt"
	asBytes, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	_, err = wr.Write(asBytes)
	return err
}

func GeneratePkgsList(ctx context.Context, wr io.Writer, searcher *blugesearcher.PackagesSearcher) {
	pkgsSeq := Must(searcher.SearchPackages(ctx, "", search.Opts{}))
	for pkg := range pkgsSeq {
		fmt.Fprintln(wr, pkg)
	}
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
