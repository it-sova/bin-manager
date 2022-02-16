package helpers

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// CreateDirIfNotExists creates directory by passed path if it's not exists
func CreateDirIfNotExists(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		log.Debug("Failed to stat directory %v, let's create one ", dir)
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create directory, %w", err)
		}
	}
	return nil
}
