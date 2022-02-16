package packets

import (
	"github.com/it-sova/bin-manager/helpers"
	"github.com/it-sova/bin-manager/remote"
	log "github.com/sirupsen/logrus"
)

// ListVersions parses remote to get available packet versions
func (p Packet) ListVersions() error {
	remote, err := remote.FindRemote(p.URLType)
	if err != nil {
		return err
	}

	versions, err := remote.ListPacketVersions(p.URL, p.Filenames, p.VersionRegex)
	if err != nil {
		return err
	}

	log.Debug("Found versions: %#v", versions)

	latestVersion, latestVersionURL, err := helpers.GetLastMapElement(versions)
	if err != nil {
		return err
	}

	log.Infof("Going to install latest %v version %v from %v", p.Name, latestVersion, latestVersionURL)
	return nil
}

// Install installs packet to OS
func (p Packet) Install() (string, error) {
	err := p.ListVersions()
	if err != nil {
		log.Error(err)
	}

	return "", nil
}
