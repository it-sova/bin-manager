package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// CreateDirIfNotExists creates directory by passed path if it's not exists
func CreateDirIfNotExists(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		log.Debugf("Failed to stat directory %v, let's create one ", dir)
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create directory, %w", err)
		}
	}
	return nil
}

// DownloadFile will download url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Errorf("Failed to close file: %v", err)
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	return err
}
