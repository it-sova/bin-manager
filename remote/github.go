package remote

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"

	"github.com/it-sova/bin-manager/helpers"
)

type githubRemote struct {
	name string
}

// NewGithubRemote creates new github remote
func NewGithubRemote() Remote {
	return githubRemote{
		name: "github",
	}
}

func (r githubRemote) ListPacketVersions(packetURL *url.URL) ([]string, error) {
	var result []string
	//TODO: Regexp?
	repoDetails := helpers.RemoveEmptyElementsFromStringSlice(strings.Split(packetURL.Path, "/"))

	if len(repoDetails) != 2 {
		return result, fmt.Errorf("Failed to get user and repo from packet URL %#v", repoDetails)
	}

	ctx := context.Background()
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(ctx, repoDetails[0], repoDetails[1], &github.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	for _, release := range releases {
		log.Info(*release.TagName)
		log.Info(*release.TarballURL)
		for _, asset := range release.Assets {
			log.Info(*asset.Name)
			log.Info(*asset.BrowserDownloadURL)
		}

	}
	return result, nil
}

func (r githubRemote) GetName() string {
	return r.name
}
