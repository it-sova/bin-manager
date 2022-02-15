package packets

import (
	"net/url"

	"gopkg.in/yaml.v3"
)

type Packet struct {
	Name        string
	URL         url.URL `yaml:"-"`
	URLRaw      string  `yaml:"url"`
	UrlType     string  `yaml:"urlType"`
	Description string
}

// TODO: Custom Unmarshaller for URL?
func NewPacket(config []byte) (Packet, error) {
	packet := Packet{}
	err := yaml.Unmarshal(config, &packet)
	if err != nil {
		return packet, err
	}

	url, err := url.Parse(packet.URLRaw)
	if err != nil {
		return packet, err
	}

	packet.URL = *url

	return packet, nil
}
