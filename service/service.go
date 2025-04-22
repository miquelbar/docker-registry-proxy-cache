package service

import (
	"docker-registry-proxy-cache/client"
	"fmt"
	"strings"
)

func GetRegistryClient(repo string) (client.RegistryClient, error) {
	// select registry from repo name prefix
	switch {
	case strings.HasPrefix(repo, "docker.io"), !strings.Contains(repo, "."):
		return client.NewDockerClient(), nil
	default:
		return nil, fmt.Errorf("unsupported registry for repo: %s", repo)
	}
}
