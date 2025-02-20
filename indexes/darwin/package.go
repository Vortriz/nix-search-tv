package darwin

import "github.com/3timeslazy/nix-search-tv/indexer"

type Package struct {
	indexer.Package
	Type        string   `json:"type"`
	Default     string   `json:"default"`
	Example     string   `json:"example"`
	DeclaredBy  []string `json:"declarations"`
	Description string   `json:"description"`
}
