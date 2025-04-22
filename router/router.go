package router

import (
	handler "docker-registry-proxy-cache/handlers"
	v2handler "docker-registry-proxy-cache/handlers/v2"
	"net/http"
	"os"

	"docker-registry-proxy-cache/storage"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/ping", handler.PingHandler)

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	store := storage.NewLocalBlobStorage(path + "/tmp/") // TODO:  select depending on the environment (local or S3)
	v2handler.RegisterV2Routes(r.Group("/v2"), store)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
	})

	return r
}
