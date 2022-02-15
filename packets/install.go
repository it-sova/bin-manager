package packets

import (
	"github.com/it-sova/bin-manager/remote"
	log "github.com/sirupsen/logrus"
)

func (p packet) ListVersions() error {
	remote, err := remote.FindRemote(p.UrlType)
	if err != nil {
		return err
	}

	versions, err := remote.ListPacketVersions(p.Url)
	if err != nil {
		return err
	}

	log.Infof("Found versions: %+v", versions)
	return nil
}

func (p packet) Install() (string, error) {
	err := p.ListVersions()
	if err != nil {
		log.Error(err)
	}

	return "", nil
}
