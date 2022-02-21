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

func (p *Packet) LatestVersion() (Version, error) {
	if len(p.Versions) > 0 {
		return p.Versions[0], nil
	}

	return Version{}, fmt.Errorf("version list has 0 elements")
}

// Install installs packet to OS
func (p *Packet) Install(installPath, version string) error {

	var installVersion Version
	p.FetchVersions()
	if version != "" {
		// TODO: Find version, use it
	} else {
		packetVersion, err := p.LatestVersion()
		if err != nil {
			return err
		}

		installVersion = packetVersion
	}

	log.Infof(
		"Going to install latest %v version %v from %v",
		p.Name,
		installVersion.Version,
		installVersion.AssetURL,
	)

	targetPath := path.Join(installPath, p.Name)
	log.Infof("Installing to %v", targetPath)

	err := helpers.DownloadFile(targetPath, installVersion.AssetURL)
	if err != nil {
		return fmt.Errorf("failed to install %v from %v: %v", p.Name, installVersion.AssetURL, err)
	}

	log.Debugf("Changing perms on %v", targetPath)
	// Set file execute permissions
	return os.Chmod(targetPath, 0755)
}
