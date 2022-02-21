package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// fileSystemRepo represents repository
type fileSystemRepo struct {
	name string
	path string
}

//NewFileSystemRepo creates new fileSystemRepo instance and validates path to packet repository
func NewFileSystemRepo(p string) (fileSystemRepo, error) {
	if p == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fileSystemRepo{}, fmt.Errorf("Failed to get current working dir, %w", err)
		}
		p = path.Join(wd, "dist")
	}

	dirInfo, err := os.Stat(p)
	if err != nil {
		return fileSystemRepo{}, fmt.Errorf("Unable to stat path %v: %w", p, err)
	}

	if !dirInfo.IsDir() {
		return fileSystemRepo{}, fmt.Errorf("Path %v is not a directory", p)
	}

	return fileSystemRepo{
		name: "FileSystemRepo",
		path: p,
	}, nil
}

// ScanPackets returns list of all packets in repository
func (r fileSystemRepo) ScanPackets() []string {
	var result []string
	files, err := ioutil.ReadDir(r.path)
	if err != nil {
		return result
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			result = append(result, file.Name())
		}
	}

	return result
}

// GetPacketConfig reads packet config from file and returns it
func (r fileSystemRepo) GetPacketConfig(packet string) ([]byte, error) {
	config, err := ioutil.ReadFile(path.Join(r.path, packet))
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to read packet config, %w", err)
	}

	return config, nil
}

// GetName getter for repo name
func (r fileSystemRepo) GetName() string {
	return r.name
}

// GetPath getter for repo path
func (r fileSystemRepo) GetPath() string {
	return r.path
}
