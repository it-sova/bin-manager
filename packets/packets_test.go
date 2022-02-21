package packets

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/suite"
)

type packetSuite struct {
	suite.Suite
	packetConfig rawPacket
	rawConfig    []byte
}

func (s *packetSuite) SetupTest() {
	var err error
	s.packetConfig = rawPacket{
		Name:         "jq",
		URL:          "https://github.com/stedolan/jq",
		URLType:      "github",
		Description:  "jq packet",
		VersionRegex: "^jq-(.+)$",
		FilenameTemplates: []string{
			"jq-{{ .os }}-{{ .arch }}",
			"jq-{{ .os }}{{ .arch }}",
			"jq-{{ .os }}-{{ .arch }}-static",
		},
	}

	s.rawConfig, err = yaml.Marshal(s.packetConfig)
	if err != nil {
		log.Errorf("Failed to setup test suite, %v", err)
	}
}

func (s *packetSuite) TestConstructorNotPanics() {
	s.Assert().NotPanics(func() {
		packet, err := New(s.rawConfig)
		s.Assert().NoError(err)
		s.Assert().Equal(s.packetConfig.Name, packet.Name)
	})
}

func (s *packetSuite) TestIncorrectConfigHandled() {
	s.Assert().NotPanics(func() {
		_, err := New([]byte{65, 66, 67, 226, 130, 172})
		s.Assert().Error(err)
		s.Assert().Contains(err.Error(), "cannot unmarshal")
	})
}

func (s *packetSuite) TestIncorrectURLHandled() {
	s.Assert().NotPanics(func() {
		s.packetConfig.URL = "http://$$!@!@#$#%^%#&"
		out, _ := yaml.Marshal(s.packetConfig)
		_, err := New(out)
		s.Assert().Error(err)
		s.Assert().Contains(err.Error(), "invalid URL")
	})
}

func (s *packetSuite) TestIncorrectRegexHandled() {
	s.Assert().NotPanics(func() {
		s.packetConfig.VersionRegex = `^\/(?!\/)(.*?)`
		out, _ := yaml.Marshal(s.packetConfig)
		_, err := New(out)
		s.Assert().Error(err)
		s.Assert().Contains(err.Error(), "error parsing regexp")
	})
}

func (s *packetSuite) TestIncorrectFilenameTemplateHandled() {
	s.Assert().NotPanics(func() {
		s.packetConfig.FilenameTemplates[0] = "{{ !.wRong }}"
		out, _ := yaml.Marshal(s.packetConfig)
		_, err := New(out)
		s.Assert().Error(err)
		s.Assert().Contains(err.Error(), "template: filename:1: unexpected")
	})
}

func TestPackets(t *testing.T) {
	suite.Run(t, new(packetSuite))
}
