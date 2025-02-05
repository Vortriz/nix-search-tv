package tv

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/3timeslazy/nix-search-tv/nixpkgs"
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

	newpkg := nixpkgs.Package{
		FullName: fullPkgName,
		Meta: nixpkgs.Meta{
			Description:     pkg.Description,
			LongDescription: pkg.LongDescription,
			MainProgram:     pkg.MainProgram,
			Homepages:       pkg.Homepages,
			Unfree:          pkg.Unfree,
			Name:            pkg.Name,
		},
		Version: pkg.Version,
	}
	for _, l := range pkg.Licenses {
		newpkg.Meta.Licenses = append(newpkg.Meta.Licenses, nixpkgs.License{
			FullName: l,
		})
	}

	nixpkgs.Preview(out, newpkg)
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
