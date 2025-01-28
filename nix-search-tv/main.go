package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"libdb.so/nix-search/search"
	"libdb.so/nix-search/search/searchers/blugesearcher"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("command required")
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
			fmt.Println("package name required")
			return
		}

		pkgPath := os.Args[2]
		// TODO: don't panic
		pkg := Must(SearchPkg(ctx, searcher, pkgPath))
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(pkg.Package)
		return
	}

	fmt.Printf("unknown command %q\n", os.Args[1])
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
		return search.SearchedPackage{}, err
	}

	for pkg := range pkgs {
		if pkg.Path == pkgPath {
			return pkg, nil
		}
	}

	return search.SearchedPackage{}, errors.New("not found")
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
