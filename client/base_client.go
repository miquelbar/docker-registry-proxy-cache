package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BaseRegistryClient struct {
	registryURL string
	authURL     string
	client      *http.Client
}

var (
	acceptedManifestTypes = map[string]struct{}{
		"application/vnd.docker.distribution.manifest.v2+json":      {},
		"application/vnd.docker.distribution.manifest.v1+prettyjws": {},
		"application/vnd.docker.distribution.manifest.list.v2+json": {},
		"application/vnd.oci.image.manifest.v1+json":                {},
		"application/vnd.oci.image.index.v1+json":                   {},
	}
	acceptHeaderValue = buildAcceptHeader()
)

func buildAcceptHeader() string {
	types := make([]string, 0, len(acceptedManifestTypes))
	for t := range acceptedManifestTypes {
		types = append(types, t)
	}

	return strings.Join(types, ",")
}

func (brc *BaseRegistryClient) fetchToken(url string) (string, error) {
	var expectedResponse struct {
		Token string `json:"token"`
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	resp, err := brc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch token: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&expectedResponse); err != nil {
		return "", err
	}

	return expectedResponse.Token, nil
}

func (brc *BaseRegistryClient) fetchManifest(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", acceptHeaderValue)

	resp, err := brc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manifest: %s", resp.Status)
	}

	// Limit the size of the response body to 10MB to prevent excessive memory usage
	limitedReader := io.LimitReader(resp.Body, MaxManifestSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	// Inspect content type and decode accordingly

	responseHeader := resp.Header.Get("Content-Type")
	_, ok := acceptedManifestTypes[responseHeader]
	if !ok {
		return nil, fmt.Errorf("unsupported manifest type: %s", resp.Header.Get("Content-Type"))
	}
	return body, err
}

func (brc *BaseRegistryClient) fetchBlob(url, token string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := brc.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker: failed to fetch blob: %s", resp.Status)
	}

	return resp.Body, nil
}
