package darwin

import (
	"fmt"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexer"
)

type Package struct {
	indexer.Package
	Type        string   `json:"type"`
	Default     string   `json:"default"`
	Example     string   `json:"example"`
	DeclaredBy  []string `json:"declarations"`
	Description string   `json:"description"`
}

func (pkg *Package) GetSource() string {
	if len(pkg.DeclaredBy) == 1 {
		return pkg.DeclaredBy[0]
	}

	return fmt.Sprintf(
		"https://daiderd.com/nix-darwin/manual/index.html#opt-%s",

		// There are packages with quotes in their names, like
		// system.defaults.".GlobalPreferences"."com.apple.mouse.scaling". For these,
		// the quotes must replaces with "_"
		strings.ReplaceAll(pkg.Name, `"`, "_"),
	)
}

func (pkg *Package) GetHomepage() string {
	return pkg.GetSource()
}
