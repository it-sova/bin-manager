package remote

import (
	"fmt"
	"net/url"
)

// Remote represents generic remote interfaces
type Remote interface {
	GetName() string
	GetPacketAssets(*url.URL) (map[string][]string, error)
}

// List returns list of all registered remotes
func List() []Remote {
	return []Remote{
		NewGithubRemote(),
	}
}

// FindRemote finds remote by its name
func FindRemote(name string) (Remote, error) {
	for _, remote := range List() {
		if remote.GetName() == name {
			return remote, nil
		}
	}

	return nil, fmt.Errorf("failed to find %v remote", name)
}
