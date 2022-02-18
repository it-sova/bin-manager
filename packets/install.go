package packets

import (
	"fmt"
	"github.com/it-sova/bin-manager/helpers"
	"github.com/it-sova/bin-manager/remote"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func (p Packet) GetNormalizedReleases(releases map[string][]string) map[string]string {
	result := map[string]string{}
	// Format packet version, get all releases matching OS and arch
	for release, assets := range releases {
		var version string
		log.Debugf("Found packet raw version: %v", release)
		matches := p.VersionRegex.FindStringSubmatch(release)

		if len(matches) == 2 {
			version = matches[1]
		} else {
			log.Errorf("Failed to parse packet version %v", release)
			continue
		}

		log.Debugf("Packet version: %v", version)
		// Check if asset filename is matching current OS and arch
		for _, url := range assets {
			assetName := path.Base(url)
			if helpers.StringSliceHasElement(p.Filenames, assetName) {
				log.Debugf("Found matching asset: %v - %v", assetName, url)
				result[version] = url
			}
		}
	}

	return result
}

// GetPacketVersions return map of packet version and asset URL for current OS and arch
func (p Packet) GetPacketVersions() (map[string]string, error) {
	var packetVersions map[string]string
	r, err := remote.FindRemote(p.URLType)
	if err != nil {
		return map[string]string{}, err
	}

	releases, err := r.GetPacketAssets(p.URL)
	if err != nil {
		return map[string]string{}, err
	}

	packetVersions = p.GetNormalizedReleases(releases)

	log.Debug("Found versions: %#v", packetVersions)

	return packetVersions, nil
}

// ListVersions parses remote to get available packet versions
func (p Packet) ListVersions() {
	versions, err := p.GetPacketVersions()
	if err != nil {
		log.Error("Failed to list available versions")
	}
	for version, asset := range versions {
		log.Infof("Found %v %v - %v", p.Name, version, asset)
	}

}

// GetLastVersion returns latest available packet version
func (p Packet) GetLastVersion() (string, string, error) {
	versions, err := p.GetPacketVersions()
	if err != nil {
		return "", "", fmt.Errorf("Failed to list available versions")
	}

	return helpers.GetLastMapElement(versions)
}

// Install installs packet to OS
func (p Packet) Install(installPath string) error {
	//TODO: If installVersion not passed....

	latestVersion, latestVersionURL, err := p.GetLastVersion()
	if err != nil {
		return fmt.Errorf("Failed to get latest packet version: %w", err)
	}
	log.Infof("Going to install latest %v version %v from %v", p.Name, latestVersion, latestVersionURL)

	targetPath := path.Join(installPath, p.Name)
	log.Infof("Installing to %v", targetPath)

	err = helpers.DownloadFile(targetPath, latestVersionURL)
	if err != nil {
		log.Errorf("Failed to install %v from %v: %v", p.Name, latestVersionURL, err)
	}

	log.Debugf("Changing perms on %v", targetPath)
	// Set file execute permissions
	err = os.Chmod(targetPath, 0755)
	if err != nil {
		log.Errorf("Failed to change file permissions: %v", err)
	}
	return nil
}
