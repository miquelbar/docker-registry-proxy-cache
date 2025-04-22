package v2_manifest

import (
	"crypto/sha256"
	"docker-registry-proxy-cache/models"
	"docker-registry-proxy-cache/storage"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Storage storage.BlobStorage
}

func NewManifestHandler(store storage.BlobStorage) *Handler {
	return &Handler{Storage: store}
}

// HeadManifest: HEAD /v2/:name/manifests/:reference
func (h *Handler) V2HeadManifests(c *gin.Context, repo, reference string) {
	manifest, err := h.Storage.LoadManifest(repo, reference)
	if err != nil {
		log.Printf("Error loading manifest: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "manifest not found"})
		return
	}

	digest := computeDigest(manifest)
	log.Printf("Manifest digest: %s", digest)

	c.Header("Content-Type", "application/vnd.oci.image.index.v1+json")
	c.Header("Docker-Content-Digest", digest)
	c.Header("docker-distribution-api-version", "registry/2.0")
	c.Header("Content-Length", fmt.Sprintf("%d", len(manifest)))
	c.Status(http.StatusOK)
}

// GetManifest: GET /v2/:name/manifests/:reference
func (h *Handler) V2GetManifest(c *gin.Context, repo, reference string) {
	manifest, err := h.Storage.LoadManifest(repo, reference)
	if err != nil {
		log.Printf("Error loading manifest: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "manifest not found"})
		return
	}

	digest := computeDigest(manifest)
	log.Printf("Manifest digest: %s", digest)
	var mediaType models.MediaTypeHolder
	if err := json.Unmarshal(manifest, &mediaType); err != nil {
		c.Data(http.StatusNotFound, "application/json", []byte(`{"error": "manifest does not contain a valid media type"}`))
		return
	}

	c.Header("Content-Type", *mediaType.MediaType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(manifest)))
	c.Header("docker-content-digest", digest)
	c.Header("docker-distribution-api-version", "registry/2.0")

	c.Data(http.StatusOK, "application/json", manifest)
}

func computeDigest(data []byte) string {
	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:])
}
