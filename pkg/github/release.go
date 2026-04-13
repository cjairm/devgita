package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FetchLatestRelease queries the GitHub API to get the latest release version
// for the given owner/repo. Returns the version string without the "v" prefix
// (e.g. "0.44.1").
func FetchLatestRelease(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url) // #nosec G107 - URL is constructed from known constants
	if err != nil {
		return "", fmt.Errorf("failed to query GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read GitHub API response: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to parse GitHub API response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no tag_name in GitHub API response")
	}

	version := release.TagName
	if len(version) > 0 && version[0] == 'v' {
		version = version[1:]
	}
	return version, nil
}
