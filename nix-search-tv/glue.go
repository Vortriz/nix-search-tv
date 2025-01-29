package main

import (
	"fmt"
	"strings"
)

// During nix-search and television integration, the
// following transformations happen:
//
// 1. Generate Nix packages index by running `nix-search --index`.
// That will save data in "<channel/flake>.<subpackage>?.<packageName>" format
//
// 2. Generate television source_command output with `nix-search-tv index`
// That will save data in "<subpackage>?.<packageName>" format
//
// 3. Run `nix-search-tv print` (Supposed to be called by televison)
// That will return data in "<subpackage>?.<packageName>" format (from 2.)
//
// 4. Run `nix-search-tv preview {}` (Supposed to be called by televison)
// That will accept data in "<subpackage>?.<packageName>" format (from 3.)
//
// 5. Find description with `nix-search`
// Here `nix-search` accepts data in "<packageName>" format

// SearchedPath - "<channel/flake>.<subpackage>?.<packageName>"
// FullPkgName  - "<subpackage>?.<packageName>"
// PkgName      - "<packageName>"

func CutChannel(searchedPath string) string {
	_, fullPkgName, found := strings.Cut(searchedPath, ".")
	if !found {
		panic(fmt.Sprintf("unexpected searched path %q", searchedPath))
	}
	return fullPkgName
}

func CutSubpkg(fullPkgName string) string {
	path := strings.Split(fullPkgName, ".")
	return path[len(path)-1]
}
