package repo

import (
	log "github.com/sirupsen/logrus"
)

// Repo represents repository ( sorce of packet dist. files ) interface
type Repo interface {
	GetName() string
	GetPath() string
	ScanPackets() []string
	GetPacketConfig(string) ([]byte, error)
}

// List Returns slice of loaded repositories
func List() []Repo {
	var repos []Repo

	fsRepo, err := NewFileSystemRepo("")
	if err != nil {
		log.Error("Failed to init FileSystem repo, ", err.Error())
	} else {
		repos = append(repos, fsRepo)
	}

	repos = append(repos, NewGitRepo())

	return repos
}
