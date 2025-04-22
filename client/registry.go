package client

import "io"

const MaxManifestSize = 10 * 1024 * 1024 // 10MB

type RegistryClient interface {
	GetManifest(repo, tag string) ([]byte, error)
	GetBlob(repo, reference string) (io.ReadCloser, error)
}
