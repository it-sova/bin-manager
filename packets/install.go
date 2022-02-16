package packets

import (
	"github.com/it-sova/bin-manager/remote"
	log "github.com/sirupsen/logrus"
)

// ListVersions parses remote to get available packet versions
func (p Packet) ListVersions() error {
	remote, err := remote.FindRemote(p.URLType)
	if err != nil {
		return err
	}

	versions, err := remote.ListPacketVersions(p.URL)
	if err != nil {
		return err
	}

	log.Infof("Found versions: %+v", versions)
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
