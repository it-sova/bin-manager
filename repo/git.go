package repo

// TODO: To be implemented, just a placeholder atm

type GitRepo struct {
	name string
	path string
}

func NewGitRepo() GitRepo {
	return GitRepo{
		name: "GitRepo",
		path: "https://github...",
	}
}

func (r GitRepo) ScanPackets() []string {
	return []string{}
}

func (r GitRepo) GetPacketConfig(packet string) ([]byte, error) {
	return []byte{}, nil
}

func (r GitRepo) GetName() string {
	return r.name
}

func (r GitRepo) GetPath() string {
	return r.path
}
