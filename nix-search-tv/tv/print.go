package tv

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func PrintPackages(out io.Writer) error {
	indexDir, err := defaultIndexPath()
	if err != nil {
		return err
	}

	nixpkgsFile := filepath.Join(indexDir, indexFile)
	asBytes, err := os.ReadFile(nixpkgsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("index file not found. Run `nix-search-tv index` and try again")
		}
		return fmt.Errorf("failed to read %s: %w", nixpkgsFile, err)
	}

	_, err = out.Write(asBytes)
	return err
}
