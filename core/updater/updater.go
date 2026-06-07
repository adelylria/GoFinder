package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/adelylria/GoFinder/core/version"
)

type ReleaseInfo struct {
	CurrentVersion  string
	LatestVersion   string
	ReleaseURL      string
	AssetName       string
	DownloadURL     string
	UpdateAvailable bool
}

type releaseResponse struct {
	TagName string         `json:"tag_name"`
	HTMLURL string         `json:"html_url"`
	Assets  []releaseAsset `json:"assets"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func CheckLatest(ctx context.Context, current string) (ReleaseInfo, error) {
	release, err := fetchLatestRelease(ctx)
	if err != nil {
		return ReleaseInfo{CurrentVersion: current}, err
	}

	asset := pickReleaseAsset(release.Assets)
	return ReleaseInfo{
		CurrentVersion:  current,
		LatestVersion:   release.TagName,
		ReleaseURL:      release.HTMLURL,
		AssetName:       asset.Name,
		DownloadURL:     asset.BrowserDownloadURL,
		UpdateAvailable: isNewerVersion(release.TagName, current),
	}, nil
}

func Download(ctx context.Context, downloadURL string) (string, error) {
	if strings.TrimSpace(downloadURL) == "" {
		return "", fmt.Errorf("release asset not found for %s", runtime.GOOS)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	// Create a secure, unpredictable temporary file for the download
	tmpFile, err := os.CreateTemp(os.TempDir(), "gofinder-update-*.exe")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		// Attempt to remove the temp file on error
		_ = os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func fetchLatestRelease(ctx context.Context) (releaseResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", version.RepoOwner, version.RepoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return releaseResponse{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", version.RepoName)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return releaseResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return releaseResponse{}, fmt.Errorf("github release check failed: %s", resp.Status)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return releaseResponse{}, err
	}
	return release, nil
}

func pickReleaseAsset(assets []releaseAsset) releaseAsset {
	if runtime.GOOS != "windows" {
		return releaseAsset{}
	}

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
		if strings.HasSuffix(name, ".exe") {
			return asset
		}
	}
	return releaseAsset{}
}

func isNewerVersion(latest string, current string) bool {
	current = normalizeVersion(current)
	latest = normalizeVersion(latest)
	if current == "" || current == "dev" || latest == "" {
		return false
	}
	if !isNumericVersion(current) || !isNumericVersion(latest) {
		return false
	}

	currentParts := versionParts(current)
	latestParts := versionParts(latest)
	for i := range latestParts {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}
	return false
}

func isNumericVersion(value string) bool {
	for _, chunk := range strings.Split(value, ".") {
		if chunk == "" {
			return false
		}
		if _, err := strconv.Atoi(chunk); err != nil {
			return false
		}
	}
	return true
}

func normalizeVersion(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.TrimPrefix(value, "v")
	if idx := strings.IndexAny(value, "-+"); idx >= 0 {
		value = value[:idx]
	}
	return value
}

func versionParts(value string) [3]int {
	var parts [3]int
	chunks := strings.Split(value, ".")
	for i := 0; i < len(parts) && i < len(chunks); i++ {
		parts[i], _ = strconv.Atoi(chunks[i])
	}
	return parts
}
