package service

import "io"

func DownloadBlob(imageRef, repo, reference string) (io.ReadCloser, error) {
	registryClient, err := GetRegistryClient(imageRef)

	if err != nil {
		return nil, err
	}

	return registryClient.GetBlob(repo, reference)
}
