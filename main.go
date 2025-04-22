package main

import (
	"docker-registry-proxy-cache/router"
)

func main() {

	r := router.SetupRouter()
	// TODO: port should be configurable
	r.Run("0.0.0.0:5001")
}
