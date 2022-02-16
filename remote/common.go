package remote

import (
	"fmt"
	"net/url"
	"regexp"
)

// Remote interface for remotes implementation
type Remote interface {
	GetName() string
	ListPacketVersions(*url.URL, []string, *regexp.Regexp) (map[string]string, error)
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

	return nil, fmt.Errorf("Failed to find %v remote", name)
}
