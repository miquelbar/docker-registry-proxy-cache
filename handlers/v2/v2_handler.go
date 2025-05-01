package v2

import (
	v2_blobs "docker-registry-proxy-cache/handlers/v2/blobs"
	v2_manifest "docker-registry-proxy-cache/handlers/v2/manifests"
	"docker-registry-proxy-cache/storage"
	"docker-registry-proxy-cache/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func V2RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"message": "API version 2",
	})
}

func V2Handler(store storage.BlobStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageRef := c.Param("imageref")
		name := c.Param("name")
		fullPath := c.Param("fullpath")
		method := c.Request.Method

		// Define the regular expression to capture the repository, action (manifests/blobs), and reference
		supportedPaths := []string{"manifests", "blobs"}
		prefix, suffix, verb, matched := utils.SplitBySupportedPath(fullPath, supportedPaths)
		if !matched {
			c.JSON(400, gin.H{"error": "unsupported path"})
			return
		}
		repo := name + prefix
		reference := suffix
		switch verb {
		// /v2/:name/manifests/:reference
		case "manifests":
			h := v2_manifest.NewManifestHandler(store)
			if method == "HEAD" {
				h.V2HeadManifests(c, imageRef, repo, reference)
			} else if method == "GET" {
				h.V2GetManifest(c, imageRef, repo, reference)
			} else {
				c.JSON(400, gin.H{"error": "unsupported method"})
			}
		// /v2/:name/blobs/:reference
		case "blobs":
			h := v2_blobs.NewManifestHandler(store)
			if method == "GET" {
				h.V2GetBlob(c, imageRef, repo, reference)
			} else {
				c.JSON(400, gin.H{"error": "unsupported method"})
			}
		default:
			c.JSON(400, gin.H{"error": "unsupported path"})
		}
	}
}

func RegisterV2Routes(r *gin.RouterGroup, store storage.BlobStorage) {
	r.GET("/", V2RootHandler)

	// v2Manifests.RegisterV2ManifestRoutes(r, store)
	// v2blobs.RegisterV2ManifestRoutes(r, store)
	r.Any("").GET("/:imageref/:name/*fullpath", V2Handler(store))
}
