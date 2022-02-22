package packets

import (
	"bytes"
	"html/template"
	"net/url"
	"regexp"
	"runtime"

	"github.com/hashicorp/go-version"

	"github.com/it-sova/bin-manager/helpers"
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
func New(config []byte) (newPacket Packet, err error) {
	var (
		parsedPacket rawPacket
		filenames    []string
	)

	// TODO: Custom Unmarshaller for packet?
	err = yaml.Unmarshal(config, &parsedPacket)
	if err != nil {
		return newPacket, err
	}

	packetURL, err := url.Parse(parsedPacket.URL)
	if err != nil {
		return newPacket, err
	}

	regex, err := regexp.Compile(parsedPacket.VersionRegex)
	if err != nil {
		return newPacket, err
	}

	for _, possibleFilename := range parsedPacket.FilenameTemplates {
		var assetTemplate *template.Template
		assetTemplate, err = template.New("filename").Parse(possibleFilename)

		if err != nil {
			return newPacket, err
		}

		// Fill with all possible arch abbrs
		for _, arch := range helpers.ArchReference[runtime.GOARCH] {
			buf := &bytes.Buffer{}
			err = assetTemplate.Execute(buf, map[string]string{
				"os":   runtime.GOOS,
				"arch": arch,
			})

			if err != nil {
				return newPacket, err
			}

			filenames = append(filenames, buf.String())
		}
	}

	newPacket.Filenames = filenames
	newPacket.URL = packetURL
	newPacket.Name = parsedPacket.Name
	newPacket.URLType = parsedPacket.URLType
	newPacket.Description = parsedPacket.Description
	newPacket.VersionRegex = regex
	newPacket.Versions = []Version{}

	return newPacket, nil
}
