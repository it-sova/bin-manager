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
	name   string
	client *github.Client
}

// NewGithubRemote creates new GitHub remote
func NewGithubRemote() Remote {
	return githubRemote{
		name:   "github",
		client: github.NewClient(nil),
	}
}

func (r githubRemote) GetPacketAssets(packetURL *url.URL) (map[string][]string, error) {

	result := map[string][]string{}
	//TODO: Regexp?
	repoDetails := helpers.RemoveEmptyElementsFromStringSlice(strings.Split(packetURL.Path, "/"))

	if len(repoDetails) != 2 {
		return result, fmt.Errorf("failed to get user and repo from packet URL %#v", repoDetails)
	}
	releases, _, err := r.client.Repositories.ListReleases(context.Background(), repoDetails[0], repoDetails[1], &github.ListOptions{})
	if err != nil {
		log.Error(err)
	}

	for _, release := range releases {
		// We don't need any versions without assets
		if len(release.Assets) > 0 {
			result[*release.TagName] = []string{}
			for _, asset := range release.Assets {
				result[*release.TagName] = append(result[*release.TagName], *asset.BrowserDownloadURL)
			}
		}
	}

	return result, nil
}

func (r githubRemote) GetName() string {
	return r.name
}
