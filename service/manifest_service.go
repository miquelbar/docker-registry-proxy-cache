package service

import "log"

func DownloadManifest(imageRef, repo, tag string) ([]byte, error) {

	log.Printf("Downloading manifest for %s:%s from %s", repo, tag, imageRef)
	registryClient, err := GetRegistryClient(imageRef)

	if err != nil {
		return nil, err
	}

	return registryClient.GetManifest(repo, tag)
}
