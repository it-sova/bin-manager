package packets

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/it-sova/bin-manager/helpers"
	"github.com/it-sova/bin-manager/remote"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"sort"
)

func (p *Packet) NormalizeReleases(releases map[string][]string) {
	// Format packet version, get all releases matching OS and arch
	for release, assets := range releases {
		var rawVersion string
		log.Debugf("Found packet raw version: %v", release)
		// Regex required to clean version from prefixes\affixes etc
		// like jq-v1.6 -> 1.6 etc
		matches := p.VersionRegex.FindStringSubmatch(release)
		if len(matches) == 2 {
			rawVersion = matches[1]
		} else {
			log.Errorf("Failed to parse packet version %v: %#v", release, matches)
			continue
		}

		parsedVersion, err := version.NewVersion(rawVersion)
		if err != nil {
			log.Errorf("Failed to parse version: %v", err)
		}

		log.Debugf("Packet version: %v", parsedVersion)
		// Check if asset filename is matching current OS and arch
		for _, url := range assets {
			assetName := path.Base(url)
			if helpers.StringSliceHasElement(p.Filenames, assetName) {
				log.Debugf("Found matching asset: %v - %v", assetName, url)
				p.Versions = append(p.Versions, Version{
					Version:  parsedVersion,
					AssetURL: url,
				})
			}
		}
	}

	// Sort versions slice
	sort.Slice(p.Versions, func(i, j int) bool {
		return !p.Versions[i].Version.LessThan(p.Versions[j].Version)
	})

}

func (p *Packet) FetchVersions() error {
	r, err := remote.FindRemote(p.URLType)
	if err != nil {
		return err
	}

	releases, err := r.GetPacketAssets(p.URL)
	if err != nil {
		return err
	}

	p.NormalizeReleases(releases)
	log.Debug("Found versions: %#v", p.Versions)

	return nil
}

// Install installs packet to OS
func (p *Packet) Install(installPath string) error {

	p.FetchVersions()
	//TODO: If installVersion passed - use it
	packetVersion := p.Versions[0]
	log.Infof(
		"Going to install latest %v version %v from %v",
		p.Name,
		packetVersion.Version,
		packetVersion.AssetURL,
	)

	targetPath := path.Join(installPath, p.Name)
	log.Infof("Installing to %v", targetPath)

	err := helpers.DownloadFile(targetPath, packetVersion.AssetURL)
	if err != nil {
		return fmt.Errorf("failed to install %v from %v: %v", p.Name, packetVersion.AssetURL, err)
	}

	log.Debugf("Changing perms on %v", targetPath)
	// Set file execute permissions
	return os.Chmod(targetPath, 0755)
}
