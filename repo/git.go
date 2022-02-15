package repo

type gitRepo struct {
	name string
	path string
}

func NewGitRepo() gitRepo {
	return gitRepo{
		name: "GitRepo",
		path: "https://github...",
	}
}

func (r gitRepo) ScanPackets() []string {
	return []string{}
}

func (r gitRepo) GetPacketConfig(packet string) ([]byte, error) {
	return []byte{}, nil

}

func (r gitRepo) GetName() string {
	return r.name
}

func (r gitRepo) GetPath() string {
	return r.path
}
