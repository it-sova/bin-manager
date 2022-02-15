package remote

import (
	"fmt"
	"net/url"
)

type Remote interface {
	GetName() string
	ListPacketVersions(url.URL) ([]string, error)
}

func RemoteList() []Remote {
	return []Remote{
		NewGithubRemote(),
	}
}

func FindRemote(name string) (Remote, error) {
	for _, remote := range RemoteList() {
		if remote.GetName() == name {
			return remote, nil
		}
	}

	return nil, fmt.Errorf("Failed to find %v remote", name)
}
