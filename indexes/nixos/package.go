package nixos

import (
	"fmt"

	"github.com/3timeslazy/nix-search-tv/indexer"
)

type Package struct {
	indexer.Package
	Example      Example  `json:"example"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Declarations []string `json:"declarations"`
	Default      Example  `json:"default"`
}

type Example struct {
	Text string `json:"text"`
}

func (pkg *Package) GetSource() string {
	if len(pkg.Declarations) == 1 {
		return fmt.Sprintf("https://github.com/NixOS/nixpkgs/blob/nixos-unstable/%s", pkg.Declarations[0])
	}

	return fmt.Sprintf(
		"https://search.nixos.org/options?"+
			"channel=unstable"+
			"&from=0&size=1"+
			"&sort=relevances&query=%[1]s"+
			// `show` automatically expands the package definion
			// on the search page. Save users a click!
			"&show=%[1]s",
		pkg.Name,
	)
}
