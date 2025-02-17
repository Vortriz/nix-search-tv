package nur

import (
	"github.com/3timeslazy/nix-search-tv/indexes/nixpkgs"
)

type Package struct {
	nixpkgs.Package
}

func (pkg *Package) GetSource() string {
	return pkg.Meta.Position
}
