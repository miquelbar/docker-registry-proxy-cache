package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type GCRIOClient struct {
	*BaseRegistryClient
}

func NewGCRIOClient() *GCRIOClient {
	return &GCRIOClient{
		BaseRegistryClient: &BaseRegistryClient{
			registryURL: "gcr.io",
			authURL:     "gcr.io",
			client:      &http.Client{},
		},
	}
}

func (gio *GCRIOClient) GetToken(repo string) (string, error) {
	// https://gcr.io/v2/token?scope=repository:distroless/base:pull&service=gcr.io
	url := fmt.Sprintf("https://%s/v2/token?scope=repository:%s/base:pull&service=gcr.io", gio.authURL, repo)
	log.Println("Docker token URL:", url)

	return gio.fetchToken(url)
}

func (gio *GCRIOClient) GetBlob(repo, digest string) (io.ReadCloser, error) {
	url := fmt.Sprintf("https://%s/v2/%s/blobs/%s", gio.registryURL, repo, digest)
	log.Println("Docker blob URL:", url)

	token, err := gio.GetToken(repo)
	if err != nil {
		return nil, err
	}

	return gio.fetchBlob(url, token)
}

func (gio *GCRIOClient) GetManifest(repo, tag string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", gio.registryURL, repo, tag)
	log.Println("Docker manifest URL:", url)

	token, err := gio.GetToken(repo)
	if err != nil {
		return nil, err
	}
	log.Printf("GCRIO token: %s", token)
	return gio.fetchManifest(url, token)
}
