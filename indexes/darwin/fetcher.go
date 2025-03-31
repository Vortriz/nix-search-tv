package darwin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/readutil"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const htmlURL = "https://nix-darwin.github.io/nix-darwin/manual/index.html"

type Fetcher struct{}

func (f *Fetcher) GetLatestRelease(ctx context.Context, _ indexer.IndexMetadata) (string, error) {
	doc, err := htmlquery.LoadURL(htmlURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch docs: %w", err)
	}

	vnode := htmlquery.FindOne(doc, `//h2[@class="subtitle"]`)
	if vnode == nil {
		return "", errors.New("no version found")
	}

	str := htmlquery.InnerText(vnode)
	_, version, _ := strings.Cut(str, " ")
	version = strings.TrimSpace(version)

	return version, nil
}

func (f *Fetcher) DownloadRelease(ctx context.Context, _ string) (io.ReadCloser, error) {
	doc, err := htmlquery.LoadURL(htmlURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch docs: %w", err)
	}

	// TODO: A better strategy would be to find the dt's sibling node
	pkgs := htmlquery.Find(doc, `//dt//span[@class="term"]`)
	pkgsContents := htmlquery.Find(doc, "//dd")

	if len(pkgs) != len(pkgsContents) {
		return nil, fmt.Errorf("len(%d) != len(%d)", len(pkgs), len(pkgsContents))
	}

	dpkgs := map[string]Package{}

	for i, pkg := range pkgs {
		dpkg := Package{}

		pkgName := htmlquery.InnerText(pkg)
		pkgName = strings.TrimSpace(pkgName)

		dpkg.Name = pkgName

		cont := pkgsContents[i]

		ems := htmlquery.Find(cont, `//span[@class="emphasis"]//em`)
		for _, em := range ems {
			emtype := htmlquery.InnerText(em)

			switch emtype {
			case "Type:":
				dpkg.Type = handleEm(em)

			case "Default:":
				p3 := em.Parent.Parent.Parent
				dpkg.Default = handleEm(em)
				if dpkg.Default != "" {
					continue
				}
				dpkg.Default = handlePreCode(p3)

			case "Example:":
				p3 := em.Parent.Parent.Parent
				dpkg.Example = handleEm(em)
				if dpkg.Example != "" {
					continue
				}
				dpkg.Example = handlePreCode(p3)

			case "Declared by:":
				p3 := em.Parent.Parent.Parent
				refs := htmlquery.Find(p3, `//table//tr//td//code//a[@class="filename"]/@href`)
				for _, ref := range refs {
					dpkg.DeclaredBy = append(dpkg.DeclaredBy, htmlquery.InnerText(ref))
				}

				tb := htmlquery.FindOne(p3, "//table")
				p3.RemoveChild(tb)
				em.Parent.RemoveChild(em)
			}
		}

		desc := htmlquery.InnerText(cont)
		desc = strings.TrimSpace(desc)
		dpkg.Description = desc

		dpkgs[pkgName] = dpkg
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(dpkgs)
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}

	return readutil.PackagesWrapper(io.NopCloser(buf)), nil
}

func handleEm(em *html.Node) string {
	p := em.Parent
	p.RemoveChild(em)

	emvalue := htmlquery.InnerText(p.Parent)
	emvalue = strings.TrimSpace(emvalue)

	del := p.Parent
	lp := del.Parent
	lp.RemoveChild(del)

	return emvalue
}

func handlePreCode(p3 *html.Node) string {
	code := htmlquery.FindOne(p3, `//pre//code`)
	if code == nil {
		return ""
	}

	text := htmlquery.InnerText(code)
	code.Parent.RemoveChild(code)
	return text
}
