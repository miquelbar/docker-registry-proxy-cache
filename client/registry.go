package client

import (
	"io"
)

const MaxManifestSize = 10 * 1024 * 1024 // 10MB

type RegistryClient interface {
	GetToken(repo string) (string, error)
	GetManifest(repo, tag string) ([]byte, error)
	GetBlob(repo, reference string) (io.ReadCloser, error)
}
