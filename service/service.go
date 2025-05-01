package service

import (
	"docker-registry-proxy-cache/client"
	"fmt"
)

func GetRegistryClient(imageRef string) (client.RegistryClient, error) {
	// select registry from repo name prefix
	switch {
	case imageRef == "docker":
		return client.NewDockerClient(), nil
	case imageRef == "gcr.io":
		return client.NewGCRIOClient(), nil
	default:
		return nil, fmt.Errorf("unsupported registry for image reference: %s", imageRef)
	}
}
