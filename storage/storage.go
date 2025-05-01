package storage

import (
	"context"
	"io"
)

type BlobStorage interface {
	SaveBlob(ctx context.Context, fullPath string, data io.Reader) error
	LoadBlob(ctx context.Context, imageRef, name, reference string) (io.ReadCloser, error)
	LoadManifest(imageRef, name, reference string) ([]byte, error)
	StoreManifest(imageRef, name, reference string, manifest []byte) error
}
