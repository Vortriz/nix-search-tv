package cmd

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/3timeslazy/nix-search-tv/config"
	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/indices"
	"github.com/3timeslazy/nix-search-tv/indexes/optionsfile"
	"github.com/3timeslazy/nix-search-tv/indexes/renderdocs"
)

var ErrUnknownIndex = errors.New("unknown index")

func SetupIndexes(conf config.Config) ([]string, error) {
	indexNames := slices.Collect(maps.Keys(indices.BuiltinIndexes))

	for index, indexHTML := range conf.Experimental.RenderDocsIndexes {
		err := indices.Register(
			index,
			renderdocs.NewFetcher(indexHTML),
			func() indices.Pkg {
				return &renderdocs.Package{
					PageURL: indexHTML,
				}
			},
		)
		if err != nil {
			return nil, fmt.Errorf("register render_docs index %q: %w", index, err)
		}

		indexNames = append(indexNames, index)
	}

	for index, path := range conf.Experimental.OptionsFile {
		err := indices.Register(
			index,
			optionsfile.NewFetcher(path),
			func() indices.Pkg {
				return &optionsfile.Package{}
			},
		)
		if err != nil {
			return nil, fmt.Errorf("register options_file index %q: %w", index, err)
		}

		indexNames = append(indexNames, index)
	}

	return indexNames, nil
}

func GetIndexes(cacheDir string, indexNames []string) ([]indexer.Index, error) {
	indexes := []indexer.Index{}
	for _, indexName := range indexNames {
		fetcher, ok := indices.GetFetcher(indexName)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownIndex, indexName)
		}

		md, err := indexer.GetIndexMetadata(cacheDir, indexName)
		if err != nil {
			return nil, fmt.Errorf("get metadata for %q: %w", indexName, err)
		}

		indexes = append(indexes, indexer.Index{
			Name:     indexName,
			Fetcher:  fetcher,
			Metadata: md,
		})
	}

	return indexes, nil
}

// The two functions below connect the print and preview
// commands. Their logic is simple, so the only reason
// these functions exist is to keep prefix logic in one place
//
// Also, pay attention to the fact the `addIndexPrefix` puts a
// space between index and package names, while `cutIndexPrefix`
// cut without the space. That's done this way because when there is
// a space, tv and fzf consider it as two arguments. Those two arguments
// later passed by tv/fzf into the preview command as "{1}{2}" and that's
// where the space dissappears

func addIndexPrefix(index, pkg string) string {
	return index + "/ " + pkg
}

func cutIndexPrefix(s string) (string, string, bool) {
	index, pkg, ok := strings.Cut(s, "/")
	return strings.TrimSpace(index), strings.TrimSpace(pkg), ok
}
