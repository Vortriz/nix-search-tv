package homemanager

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/readutil"
)

type NixBuilder struct{}

func (NixBuilder) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	cmd := exec.Command(
		"nix", "build",
		"github:nix-community/home-manager/master#docs-json",
		"--no-write-lock-file", "--no-link", "--print-out-paths",
	)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("build options: %w", err)
	}

	path := strings.TrimSpace(string(out))
	return filepath.Join(path, "/share/doc/home-manager/options.json"), nil
}

func (NixBuilder) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	file, err := os.OpenFile(release, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("open options file: %w", err)
	}

	return readutil.PackagesWrapper(file), nil
}
