package renderdocs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/readutil"
	"github.com/3timeslazy/nix-search-tv/pkgs/renderdocs"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Package struct {
	indexer.Package
	PageURL     string   `json:"-"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Default     string   `json:"default"`
	Example     string   `json:"example"`
	DeclaredBy  []string `json:"declared_by"`
}

type Fetcher struct {
	url string
}

func NewFetcher(url string) *Fetcher {
	return &Fetcher{
		url: url,
	}
}

func (f *Fetcher) GetLatestRelease(_ context.Context, _ indexer.IndexMetadata) (string, error) {
	return time.Now().String(), nil
}

func (f *Fetcher) DownloadRelease(_ context.Context, release string) (io.ReadCloser, error) {
	var doc *html.Node
	var err error

	_, path, ok := strings.Cut(release, "file://")
	if ok {
		doc, err = htmlquery.LoadDoc(path)
	} else {
		doc, err = htmlquery.LoadURL(f.url)
	}
	if err != nil {
		return nil, fmt.Errorf("download options.xhtml: %w", err)
	}

	htmlPkgs, err := renderdocs.Parse(doc)
	if err != nil {
		return nil, fmt.Errorf("parse options.xhtml: %w", err)
	}

	pkgs := map[string]Package{}
	for name, htmlPkg := range htmlPkgs {
		pkg := Package{
			Package: indexer.Package{
				Name: htmlPkg.Name,
			},
			Example:     htmlPkg.Example,
			Type:        htmlPkg.Type,
			Description: htmlPkg.Description,
			Default:     htmlPkg.Default,
			DeclaredBy:  htmlPkg.DeclaredBy,
		}

		pkgs[name] = pkg
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(pkgs)
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}

	return readutil.PackagesWrapper(io.NopCloser(buf)), nil
}
