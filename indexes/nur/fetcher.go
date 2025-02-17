package nur

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/3timeslazy/nix-search-tv/indexer"
	"github.com/3timeslazy/nix-search-tv/indexes/readutil"
)

type Fetcher struct{}

const commitsURL = "https://api.github.com/repos/nix-community/nur-search/commits?page=1&per_page=1"

func (f *Fetcher) GetLatestRelease(ctx context.Context, md indexer.IndexMetadata) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, commitsURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("github request failed: %w", err)
	}
	defer resp.Body.Close()

	commits := []struct {
		Sha string `json:"sha"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal github response: %w", err)
	}

	if len(commits) < 1 {
		return "", fmt.Errorf("unexpected result from github: %w", err)
	}
	return commits[0].Sha, nil
}

const packagesURL = "https://raw.githubusercontent.com/nix-community/nur-search/%s/data/packages.json"

func (f *Fetcher) DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	apiurl := fmt.Sprintf(packagesURL, release)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiurl, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}

	return readutil.PackagesWrapper(resp.Body), nil
}
