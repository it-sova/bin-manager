package repo

import (
	log "github.com/sirupsen/logrus"
)

type Repo interface {
	GetName() string
	GetPath() string
	ScanPackets() []string
	GetPacketConfig(string) ([]byte, error)
}

func RepoList() []Repo {
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
