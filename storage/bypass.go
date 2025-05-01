package storage

import (
	"context"
	"docker-registry-proxy-cache/service"
	"io"
)

type ByPassStorage struct{}

func NewByPassStorage() *ByPassStorage {
	return &ByPassStorage{}
}

func (b *ByPassStorage) SaveBlob(ctx context.Context, fullPath string, data io.Reader) error {
	return nil
}
func (b *ByPassStorage) LoadBlob(ctx context.Context, imageRef, name, reference string) (io.ReadCloser, error) {
	return service.DownloadBlob(imageRef, name, reference)
}
func (b *ByPassStorage) LoadManifest(imageRef, name, reference string) ([]byte, error) {
	return service.DownloadManifest(imageRef, name, reference)
}
func (b *ByPassStorage) StoreManifest(imageRef, name, reference string, manifest []byte) error {
	return nil
}
