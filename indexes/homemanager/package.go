package homemanager

import (
	"github.com/3timeslazy/nix-search-tv/indexer"
)

type Package struct {
	indexer.Package
	Example      Example        `json:"example"`
	Type         string         `json:"type"`
	Description  string         `json:"description"`
	Declarations []Declarations `json:"declarations"`
	Default      Default        `json:"default"`
}

type Example struct {
	Text string `json:"text"`
}

type Declarations struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Default struct {
	Text string `json:"text"`
}
