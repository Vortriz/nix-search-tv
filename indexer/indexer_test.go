package indexer

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestInjectKey(t *testing.T) {
	pkg := []byte(`{ "version": "v1.0.0" }`)
	pkg = injectKey("nix-search-tv", pkg)
	assert.Equal(t, []byte(`{"_key":"nix-search-tv", "version": "v1.0.0" }`), pkg)
}
