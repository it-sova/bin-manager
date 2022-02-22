package repo

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// fileSystemRepo represents repository
type FileSystemRepo struct {
	name string
	path string
}

// NewFileSystemRepo creates new fileSystemRepo instance and validates path to packet repository
func NewFileSystemRepo(p string) (FileSystemRepo, error) {
	if p == "" {
		wd, err := os.Getwd()
		if err != nil {
			return FileSystemRepo{}, fmt.Errorf("failed to get current working dir, %w", err)
		}

		p = path.Join(wd, "dist")
	}

	dirInfo, err := os.Stat(p)
	if err != nil {
		return FileSystemRepo{}, fmt.Errorf("unable to stat path %v: %w", p, err)
	}

	if !dirInfo.IsDir() {
		return FileSystemRepo{}, fmt.Errorf("path %v is not a directory", p)
	}

	return FileSystemRepo{
		name: "FileSystemRepo",
		path: p,
	}, nil
}

// ScanPackets returns list of all packets in repository
func (r FileSystemRepo) ScanPackets() []string {
	var result []string

	files, err := os.ReadDir(r.path)

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
func (r FileSystemRepo) GetPacketConfig(packet string) ([]byte, error) {
	config, err := os.ReadFile(path.Join(r.path, packet))

	if err != nil {
		return []byte{}, fmt.Errorf("failed to read packet config, %w", err)
	}

	return config, nil
}

// GetName getter for repo name
func (r FileSystemRepo) GetName() string {
	return r.name
}

// GetPath getter for repo path
func (r FileSystemRepo) GetPath() string {
	return r.path
}
