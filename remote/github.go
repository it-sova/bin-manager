package remote

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/it-sova/bin-manager/helpers"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type githubRemote struct {
	name   string
	client *github.Client
}

// NewGithubRemote creates new GitHub remote
func NewGithubRemote() Remote {
	var client *http.Client
	token := os.Getenv("GITHUB_TOKEN")

	if len(token) > 0 {
		log.Debugf("GitHub API Token found, using token-based auth")
		client = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		)
	}

	return githubRemote{
		name:   "github",
		client: github.NewClient(client),
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
		} else {
			log.Debugf("Empty release: %v", *release.TagName)
		}
	}

	return result, nil
}

func (r githubRemote) GetName() string {
	return r.name
}
