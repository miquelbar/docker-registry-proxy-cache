package service

func DownloadManifest(repo, tag string) ([]byte, error) {

	registryClient, err := GetRegistryClient(repo)

	if err != nil {
		return nil, err
	}

	return registryClient.GetManifest(repo, tag)
}
