package packets

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/it-sova/bin-manager/helpers"
	"github.com/it-sova/bin-manager/remote"
	log "github.com/sirupsen/logrus"
)

// NormalizeReleases parses general version from raw one got from remote
// Filters assets to keep only matching current OS and arch
// Sorts resulting slice and applies it to structure
func (p *Packet) NormalizeReleases(releases map[string][]string) {
	// Format packet version, get all releases matching OS and arch
	for release, assets := range releases {
		var rawVersion string

		log.Debugf("Found packet raw version: %v", release)
		// Regex required to clean version from prefixes\affixes etc
		// like jq-v1.6 -> 1.6 etc
		matches := p.VersionRegex.FindStringSubmatch(release)
		if len(matches) == releaseRegexExpectedMatchCount {
			rawVersion = matches[1]
		} else {
			log.Errorf("Failed to parse packet version %v: %#v", release, matches)
			continue
		}

		parsedVersion, err := version.NewVersion(rawVersion)
		if err != nil {
			log.Errorf("Failed to parse version: %v", err)
			continue
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

// FetchVersions fetches versions from remote
func (p *Packet) FetchVersions() error {
	// To not fetch versions more than 1 time
	if len(p.Versions) > 0 {
		return nil
	}

	r, err := remote.FindRemote(p.URLType)
	if err != nil {
		return err
	}

	releases, err := r.GetPacketAssets(p.URL)
	if err != nil {
		return err
	}

	p.NormalizeReleases(releases)

	return nil
}

// FindVersion searches for given version of packet
func (p *Packet) FindVersion(targetVersion string) (Version, bool) {
	err := p.FetchVersions()
	if err != nil {
		return Version{}, false
	}

	for _, packetVersion := range p.Versions {
		if packetVersion.Version.String() == targetVersion {
			return packetVersion, true
		}
	}

	return Version{}, false
}

// LatestVersion returns latest packet version if possible
func (p *Packet) LatestVersion() (Version, error) {
	if err := p.FetchVersions(); err != nil {
		return Version{}, fmt.Errorf("failed to fetch packet versions, %w", err)
	}

	if len(p.Versions) > 0 {
		return p.Versions[0], nil
	}

	return Version{}, fmt.Errorf("version list has 0 elements")
}

// Install installs packet
func (p *Packet) Install(installPath, packetVersion string) error {
	err := p.FetchVersions()
	if err != nil {
		return fmt.Errorf("failed to fetch packet versions, %w", err)
	}

	installVersion, _ := p.FindVersion(packetVersion)

	log.Infof(
		"Going to install %v %v from %v",
		p.Name,
		installVersion.Version,
		installVersion.AssetURL,
	)

	targetPath := path.Join(installPath, p.Name)
	log.Infof("Installing to %v", targetPath)

	err = helpers.DownloadFile(targetPath, installVersion.AssetURL)
	if err != nil {
		return fmt.Errorf("failed to install %v from %v: %v", p.Name, installVersion.AssetURL, err)
	}

	log.Debugf("Changing perms on %v", targetPath)
	// Set file execute permissions
	return os.Chmod(targetPath, helpers.FileExecutableChmod)
}
