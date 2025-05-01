package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type DockerClient struct {
	*BaseRegistryClient
}

func NewDockerClient() *DockerClient {
	return &DockerClient{
		BaseRegistryClient: &BaseRegistryClient{
			registryURL: "registry-1.docker.io",
			authURL:     "auth.docker.io",
			client:      &http.Client{},
		},
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
	pullScope := dc.GetRepoName(repo) + ":pull"
	url := fmt.Sprintf("https://%s/token?service=registry.docker.io&scope=repository:%s", dc.authURL, pullScope)
	log.Println("Docker token URL:", url)

	return dc.fetchToken(url)
}

func (dc *DockerClient) GetBlob(repo, digest string) (io.ReadCloser, error) {

	url := fmt.Sprintf("https://%s/v2/%s/blobs/%s", dc.registryURL, dc.GetRepoName(repo), digest)
	log.Println("gcr.io blob URL:", url)

	token, err := dc.GetToken(repo)
	if err != nil {
		return nil, err
	}

	return dc.fetchBlob(url, token)
}

// GetManifest retrieves the manifest of a Docker image from the registry.
// It takes the repository name and tag as input and returns the manifest as a byte slice or an error.
func (dc *DockerClient) GetManifest(repo, tag string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", dc.registryURL, dc.GetRepoName(repo), tag)
	log.Println("gcr.io manifest URL:", url)

	token, err := dc.GetToken(repo)
	if err != nil {
		return nil, err
	}
	return dc.fetchManifest(url, token)
}
