// Package indices is a wrapper around indexes/ package providing
// a more convenient API for validating available indexes and picking up
// their fetches and previews
//
// Also, as a lot of code uses `indexes` word for variable names, this package
// creates an `indices` alias to avoid conflicts with `indexes`
package indices

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/homemanager"
	"github.com/3timeslazy/nix-search-tv/indexes/nixpkgs"
	"github.com/3timeslazy/nix-search-tv/indexes/nur"
)

const (
	Nixpkgs     = "nixpkgs"
	HomeManager = "home-manager"
	Nur         = "nur"
)

var Indexes = map[string]bool{
	Nixpkgs:     true,
	HomeManager: true,
	Nur:         true,
}

var Fetchers = map[string]indexer.Fetcher{
	Nixpkgs:     &nixpkgs.Fetcher{},
	HomeManager: &homemanager.Fetcher{},
	Nur:         &nur.Fetcher{},
}

func Preview(out io.Writer, index string, pkg json.RawMessage) error {
	switch index {
	case Nixpkgs:
		nixpkg := nixpkgs.Package{}
		if err := json.Unmarshal(pkg, &nixpkg); err != nil {
			return fmt.Errorf("unmarshal package: %w", err)
		}
		nixpkgs.Preview(out, nixpkg)

	case HomeManager:
		hmpkg := homemanager.Package{}
		if err := json.Unmarshal(pkg, &hmpkg); err != nil {
			return fmt.Errorf("unmarshal package: %w", err)
		}
		homemanager.Preview(out, hmpkg)

	case Nur:
		nurpkg := nur.Package{}
		if err := json.Unmarshal(pkg, &nurpkg); err != nil {
			return fmt.Errorf("unmarshal package: %w", err)
		}
		nur.Preview(out, nurpkg)

	default:
		return errors.New("unknown index")
	}

	return nil
}

func SourcePreview(out io.Writer, index string, pkg json.RawMessage) error {
	var src interface {
		GetSource() string
	}

	switch index {
	case Nixpkgs:
		src = &nixpkgs.Package{}

	case HomeManager:
		src = &homemanager.Package{}

	case Nur:
		src = &nur.Package{}

	default:
		return errors.New("unknown index")
	}

	if err := json.Unmarshal(pkg, &src); err != nil {
		return fmt.Errorf("unmarshal package: %w", err)
	}

	_, err := out.Write([]byte(src.GetSource()))
	return err
}
