package storage

import (
	"context"
	"docker-registry-proxy-cache/service"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type LocalBlobStorage struct {
	BasePath string
}

func NewLocalBlobStorage(basePath string) *LocalBlobStorage {
	return &LocalBlobStorage{BasePath: basePath}
}

func (l *LocalBlobStorage) GetBlobPath(imageRef, name, reference string) string {
	// Construct the path to the manifest file
	return filepath.Join(l.BasePath, imageRef, name, "blobs", "rootfs", reference)
}

func (l *LocalBlobStorage) SaveBlob(ctx context.Context, fullPath string, data io.Reader) error {
	// Ensure parent directories exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, data)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

func (l *LocalBlobStorage) LoadBlob(ctx context.Context, imageRef, name, reference string) (io.ReadCloser, error) {
	fullPath := l.GetBlobPath(imageRef, name, reference)
	log.Printf("Attempting to load blob: name=%s reference=%s path=%s", name, reference, fullPath)

	file, err := os.Open(fullPath)
	if err == nil {
		log.Printf("Blob found locally at %s", fullPath)
		return file, nil
	}
	log.Printf("Blob not found locally. Attempting download: name=%s reference=%s", name, reference)
	remoteBlob, err := service.DownloadBlob(imageRef, name, reference)
	if err != nil {
		return nil, fmt.Errorf("could not download blob (name=%s, reference=%s): %w", name, reference, err)
	}
	// SaveBlob should consume the stream
	if err := l.SaveBlob(ctx, fullPath, remoteBlob); err != nil {
		log.Printf("Failed to save downloaded blob: %v", err)
	}
	defer remoteBlob.Close()

	// Now re-open to return a fresh ReadCloser
	file, err = os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to re-open blob after saving: %w", err)
	}
	return file, nil
}

func (l *LocalBlobStorage) GetManifestPath(imageRef, name, reference string) string {
	// Construct the path to the manifest file
	var lastPath string
	// TODO: Handle other digest types
	if strings.HasPrefix(reference, "sha256:") {
		lastPath = filepath.Join("blobs", "sha256", reference[7:]) // Remove the "sha256:" prefix
	} else {
		lastPath = filepath.Join(reference + ".json")
	}

	return filepath.Join(l.BasePath, imageRef, name, lastPath)
}

func (l *LocalBlobStorage) StoreManifest(imageRef, name, reference string, manifest []byte) error {
	// Save the manifest to local storage
	key := l.GetManifestPath(imageRef, name, reference)
	log.Print("Saving manifest to ", key)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(key), 0755); err != nil {
		return err
	}
	// Write the JSON data to a file
	file, err := os.OpenFile(key, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(manifest)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalBlobStorage) LoadManifest(imageRef, name, reference string) ([]byte, error) {
	// Save the manifest to local storage
	key := l.GetManifestPath(imageRef, name, reference)
	log.Print("Loading manifest from ", key)

	var manifest []byte

	file, err := os.Open(key)
	if err != nil {
		log.Printf("Manifest for %s/%s not found, downloading manifest...", name, reference)
		manifest, err = service.DownloadManifest(imageRef, name, reference)
		if err != nil {
			return nil, fmt.Errorf("failed to download manifest: %w", err)
		}
		// Save the manifest to local storage
		if err := l.StoreManifest(imageRef, name, reference, manifest); err != nil {
			return nil, fmt.Errorf("failed to store manifest: %w", err)
		}
	} else {
		// Decode the JSON data
		log.Printf("Loading manifest from %s", key)
		manifest, err = io.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	return manifest, nil
}
