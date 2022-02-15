package remote

import (
	"context"

	"github.com/google/go-github/v42/github"
	log "github.com/sirupsen/logrus"
)

type githubRemote struct {
	name string
}

func NewGithubRemote() Remote {
	return githubRemote{
		name: "github",
	}
}

func (r githubRemote) ListPacketVersions(string) ([]string, error) {
	client := github.NewClient(nil)

	ctx := context.Background()

	releases, _, err := client.Repositories.ListReleases(ctx, "stedolan", "jq", &github.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	for _, release := range releases {
		log.Info(*release.TagName)
		log.Info(*release.TarballURL)

	}
	return []string{}, nil
}

func (r githubRemote) GetName() string {
	return r.name
}
