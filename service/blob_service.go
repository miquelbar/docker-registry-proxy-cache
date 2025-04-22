package service

import "io"

func DownloadBlob(repo, reference string) (io.ReadCloser, error) {
	registryClient, err := GetRegistryClient(repo)

	if err != nil {
		return nil, err
	}

	return registryClient.GetBlob(repo, reference)
}
