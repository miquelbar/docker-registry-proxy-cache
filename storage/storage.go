package storage

import (
	"context"
	"io"
)

type BlobStorage interface {
	SaveBlob(ctx context.Context, key string, data io.Reader) error
	LoadBlob(ctx context.Context, name, reference string) (io.ReadCloser, error)
	LoadManifest(name, reference string) ([]byte, error)
	StoreManifest(name, reference string, manifest []byte) error
	GetManifestPath(name, reference string) string
}
