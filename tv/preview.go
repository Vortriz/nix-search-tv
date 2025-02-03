package tv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/3timeslazy/nix-search-tv/style"
	"libdb.so/nix-search/search"
	"libdb.so/nix-search/search/searchers/blugesearcher"
)

func PrintPreview(ctx context.Context, out io.Writer, fullPkgName string) error {
	searcher, err := blugesearcher.Open("")
	if err != nil {
		return fmt.Errorf("nix-search error: %w", err)
	}
	defer searcher.Close()

	pkg, err := searchPkg(ctx, searcher, fullPkgName)
	if err != nil {
		return fmt.Errorf("search package: %w", err)
	}

	printPreview(out, pkg)
	return nil
}

func searchPkg(ctx context.Context, searcher *blugesearcher.PackagesSearcher, fullPkgName string) (search.SearchedPackage, error) {
	query := cutSubpkg(fullPkgName)
	pkgs, err := searcher.SearchPackages(ctx, query, search.Opts{
		Exact: true,
	})
	if err != nil {
		return search.SearchedPackage{}, fmt.Errorf("nix-search error: %w", err)
	}

	for pkg := range pkgs {
		if cutChannel(pkg.Path) == fullPkgName {
			return pkg, nil
		}
	}

	return search.SearchedPackage{}, errors.New("not found")
}

func printPreview(out io.Writer, pkg search.SearchedPackage) {
	styler := style.StyledText

	// title
	fmt.Fprint(out, styler.Red(styler.Bold(pkg.Name)))
	fmt.Fprint(out, " ", styler.Dim("("+pkg.Version+")"))
	if pkg.Broken {
		fmt.Fprint(out, " ", styler.Red("(broken)"))
	}
	fmt.Fprintln(out)

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
