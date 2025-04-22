package v2_blobs

import (
	"docker-registry-proxy-cache/storage"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Storage storage.BlobStorage
}

func NewManifestHandler(store storage.BlobStorage) *Handler {
	return &Handler{Storage: store}
}

func (h *Handler) V2GetBlob(c *gin.Context, repo, reference string) {
	blob, err := h.Storage.LoadBlob(c, repo, reference)
	if err != nil {
		c.JSON(404, gin.H{"message": "blob not found", "error": err.Error()})
		return
	}
	defer blob.Close()

	// Read from the blob into a byte slice
	data, err := io.ReadAll(blob)
	if err != nil {
		log.Printf("failed to read blob: %s\n", err)
		c.JSON(500, gin.H{"message": "failed to read blob", "error": err.Error()})
		return
	}

	c.Data(200, "application/octet-stream", data)
}
