package packets

import (
	"html/template"
	"net/url"
	"regexp"

	"gopkg.in/yaml.v3"
)

// rawPacket represents raw packet as it stores in .yaml manifests
type rawPacket struct {
	Name             string `yaml:"Name"`
	URL              string `yaml:"URL"`
	URLType          string `yaml:"URLType"`
	Description      string `yaml:"Description"`
	VersionRegex     string `yaml:"VersionRegex"`
	FilenameTemplate string `yaml:"FilenameTemplate"`
}

// Packet represents parsed packet
type Packet struct {
	Name             string
	URL              *url.URL
	URLType          string
	Description      string
	VersionRegex     *regexp.Regexp
	FilenameTemplate *template.Template
}

// NewPacket builds Packet struct from rawPacket
func NewPacket(config []byte) (Packet, error) {
	rawPacket := rawPacket{}
	packet := Packet{}
	// TODO: Custom Unmarshaller for URL?
	err := yaml.Unmarshal(config, &rawPacket)
	if err != nil {
		return packet, err
	}

	url, err := url.Parse(rawPacket.URL)
	if err != nil {
		return packet, err
	}

	regex, err := regexp.Compile(rawPacket.VersionRegex)
	if err != nil {
		return packet, err
	}

	template, err := template.New("filename").Parse(rawPacket.FilenameTemplate)
	if err != nil {
		return packet, err
	}

	packet.URL = url
	packet.Name = rawPacket.Name
	packet.URLType = rawPacket.URLType
	packet.Description = rawPacket.Description
	packet.VersionRegex = regex
	packet.FilenameTemplate = template

	return packet, nil
}
