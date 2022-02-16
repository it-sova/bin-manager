package remote

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
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

func (r githubRemote) ListPacketVersions(packetURL *url.URL,
	filenames []string,
	versionRegexp *regexp.Regexp) (map[string]string, error) {

	result := map[string]string{}
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

	// Format packet version, get all releases matching OS and arch
	for _, release := range releases {
		var version string
		log.Debugf("Found packet raw version: %v", *release.TagName)
		matches := versionRegexp.FindStringSubmatch(*release.TagName)

		if len(matches) == 2 {
			version = matches[1]
		} else {
			log.Errorf("Failed to parse packet version %v", *release.TagName)
			continue
		}

		log.Debugf("Packet version: %v", version)
		// Check if asset filename is matching current OS and arch
		for _, asset := range release.Assets {
			if helpers.StringSliceHasElement(filenames, *asset.Name) {
				log.Debugf("Found matching asset: %v - %v", *asset.Name, *asset.BrowserDownloadURL)
				result[version] = *asset.BrowserDownloadURL
			}
		}
	}
	return result, nil
}

func (r githubRemote) GetName() string {
	return r.name
}
