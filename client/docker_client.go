package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type DockerClient struct {
	registryURL string
	authURL     string
	client      *http.Client
}

func NewDockerClient() *DockerClient {
	return &DockerClient{
		registryURL: "registry-1.docker.io",
		authURL:     "auth.docker.io",
		client:      &http.Client{},
	}
}

// Return if I need to use library or not
func (dc *DockerClient) GetRepoName(repo string) string {
	// Check if the repository name contains a slash
	// If it does, we need to adjust the scope accordingly
	// This is a workaround for the Docker Hub API, which requires the full repository name
	// in the scope parameter
	if count := strings.Count(repo, "/"); count > 0 {
		return repo
	}
	return "library/" + repo
}

func (dc *DockerClient) GetToken(repo string) (string, error) {
	// Define the struct locally
	var expectedResponse struct {
		Token string `json:"token"`
	}

	pullScope := dc.GetRepoName(repo) + ":pull"
	url := fmt.Sprintf("https://%s/token?service=registry.docker.io&scope=repository:%s", dc.authURL, pullScope)
	log.Println("Docker token URL:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	resp, err := dc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("docker: failed to fetch token: %s", resp.Status)
	}

	// Decode JSON body
	if err := json.NewDecoder(resp.Body).Decode(&expectedResponse); err != nil {
		return "", err
	}

	return expectedResponse.Token, nil
}

func (dc *DockerClient) GetBlob(repo, digest string) (io.ReadCloser, error) {

	url := fmt.Sprintf("https://%s/v2/%s/blobs/%s", dc.registryURL, dc.GetRepoName(repo), digest)
	log.Println("Docker blob URL:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	cliToken, err := dc.GetToken(repo)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cliToken)

	resp, err := dc.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker: failed to fetch blob: %s", resp.Status)
	}

	return resp.Body, nil
}

// GetManifest retrieves the manifest of a Docker image from the registry.
// It takes the repository name and tag as input and returns the manifest as a byte slice or an error.
func (dc *DockerClient) GetManifest(repo, tag string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", dc.registryURL, dc.GetRepoName(repo), tag)
	log.Println("Docker manifest URL:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	cliToken, err := dc.GetToken(repo)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cliToken)
	req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")

	resp, err := dc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker: failed to fetch manifest for repo '%s' and tag '%s': %s", repo, tag, resp.Status)
	}

	// Limit the size of the response body to 10MB to prevent excessive memory usage
	limitedReader := io.LimitReader(resp.Body, MaxManifestSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	// Inspect content type and decode accordingly
	switch resp.Header.Get("Content-Type") {
	case "application/vnd.oci.image.index.v1+json",
		"application/vnd.oci.image.manifest.v1+json",
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.docker.distribution.manifest.v2+json":
		return body, err
	default:
		return nil, fmt.Errorf("unsupported manifest type: %s", resp.Header.Get("Content-Type"))
	}
}
