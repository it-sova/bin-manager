package packets

import (
	"bytes"
	"github.com/hashicorp/go-version"
	"html/template"
	"net/url"
	"regexp"
	"runtime"

	"github.com/it-sova/bin-manager/helpers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// rawPacket represents raw packet as it stores in .yaml manifests
type rawPacket struct {
	Name              string   `yaml:"Name"`
	URL               string   `yaml:"URL"`
	URLType           string   `yaml:"URLType"`
	Description       string   `yaml:"Description"`
	VersionRegex      string   `yaml:"VersionRegex"`
	FilenameTemplates []string `yaml:"FilenameTemplates"`
}

// Version represents packet version
type Version struct {
	Version  *version.Version
	AssetURL string
}

// Packet represents parsed packet
type Packet struct {
	Name         string
	URL          *url.URL
	URLType      string
	Description  string
	VersionRegex *regexp.Regexp
	Filenames    []string
	Versions     []Version
}

// New builds Packet struct from rawPacket
func New(config []byte) (Packet, error) {
	rawPacket := rawPacket{}
	packet := Packet{}
	// TODO: Custom Unmarshaller for packet?
	err := yaml.Unmarshal(config, &rawPacket)
	if err != nil {
		return packet, err
	}

	log.Debug("RawPacket -> %#v", rawPacket)

	url, err := url.Parse(rawPacket.URL)
	if err != nil {
		return packet, err
	}

	regex, err := regexp.Compile(rawPacket.VersionRegex)
	if err != nil {
		return packet, err
	}

	var filenames []string
	for _, possibleFilename := range rawPacket.FilenameTemplates {
		template, err := template.New("filename").Parse(possibleFilename)
		if err != nil {
			return packet, err
		}

		// Fill with all possible arch abbrs
		for _, arch := range helpers.ArchReference[runtime.GOARCH] {
			buf := &bytes.Buffer{}
			err = template.Execute(buf, map[string]string{
				"os":   runtime.GOOS,
				"arch": arch,
			})

			if err != nil {
				return packet, err
			}
			filenames = append(filenames, buf.String())
		}

	}

	packet.Filenames = filenames
	packet.URL = url
	packet.Name = rawPacket.Name
	packet.URLType = rawPacket.URLType
	packet.Description = rawPacket.Description
	packet.VersionRegex = regex
	packet.Versions = []Version{}

	log.Debug("Packet -> %#v", packet)
	return packet, nil
}
