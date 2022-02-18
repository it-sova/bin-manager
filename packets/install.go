package packets

import (
	"fmt"
	"github.com/it-sova/bin-manager/helpers"
	"github.com/it-sova/bin-manager/remote"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

// GetPacketVersions return map of packet version and asset URL for current OS and arch
func (p Packet) GetPacketVersions() (map[string]string, error) {
	r, err := remote.FindRemote(p.URLType)
	if err != nil {
		return map[string]string{}, err
	}

	versions, err := r.GetPacketAssets(p.URL, p.Filenames, p.VersionRegex)
	if err != nil {
		return map[string]string{}, err
	}

	log.Debug("Found versions: %#v", versions)

	return versions, nil
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
