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
	packet       Packet
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

	s.packet, _ = New(s.rawConfig)
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

func (s *packetSuite) TestValidRegexShouldBeParsed() {
	validRegexps := []string{
		`^v([\d{0,10}\.{0,1}]+).*$`,
		`^v(.+)$`,
		`^jq-(.+)$`,
	}
	for _, regex := range validRegexps {
		s.Assert().NotPanics(func() {
			s.packetConfig.VersionRegex = regex
			out, _ := yaml.Marshal(s.packetConfig)
			_, err := New(out)
			s.Assert().NoError(err)
		})
	}
}

func (s *packetSuite) TestValidTemplatesShouldBeParsed() {
	validTemplates := []string{
		"yq_{{ .os }}_{{ .arch }}",
		"k3d-{{ .os }}-{{ .arch }}",
		"jq-{{ .os }}-{{ .arch }}",
		"jq-{{ .os }}{{ .arch }}",
		"jq-{{ .os }}-{{ .arch }}-static",
	}

	s.Assert().NotPanics(func() {
		s.packetConfig.FilenameTemplates = validTemplates
		out, _ := yaml.Marshal(s.packetConfig)
		_, err := New(out)
		s.Assert().NoError(err)
	})
}

func (s *packetSuite) TestNormalizeReleasesNotPanics() {
	expectedVersions := []string{"1.6.0", "1.5.0"}
	releases := map[string][]string{
		"jq-1.5": {
			"https://github.com/stedolan/jq/releases/download/jq-1.4/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc2/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc1/jq-linux-x86_64-static",
		},
		"jq-1.6": {
			"https://github.com/stedolan/jq/releases/download/jq-1.4/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc2/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc1/jq-linux-x86_64-static",
		},
	}

	s.NotPanics(func() {
		s.packet.NormalizeReleases(releases)
	})

	for _, packet := range s.packet.Versions {
		s.Assert().Contains(expectedVersions, packet.Version.String())
	}
}

func (s *packetSuite) TestNormalizeReleasesFailedToParseVersion() {
	releases := map[string][]string{
		"!jq-1.5": {
			"https://github.com/stedolan/jq/releases/download/jq-1.4/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc2/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc1/jq-linux-x86_64-static",
		},
	}

	s.NotPanics(func() {
		s.packet.NormalizeReleases(releases)
		s.Assert().Len(s.packet.Versions, 0)
	})
}

func (s *packetSuite) TestNormalizeReleasesWrongVersion() {
	releases := map[string][]string{
		"jq-1.5.XXZ": {
			"https://github.com/stedolan/jq/releases/download/jq-1.4/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc2/jq-linux-x86_64",
			"https://github.com/stedolan/jq/releases/download/jq-1.5rc1/jq-linux-x86_64-static",
		},
	}

	s.NotPanics(func() {
		s.packet.NormalizeReleases(releases)
		s.Assert().Len(s.packet.Versions, 0)
	})
}

func (s *packetSuite) TestFetchVersions() {
	s.NotPanics(func() {
		err := s.packet.FetchVersions()
		s.Assert().NoError(err)
		s.Assert().NotEmpty(s.packet.Versions)
		err = s.packet.FetchVersions()
		s.Assert().NoError(err)
	})
}

func (s *packetSuite) TestFetchVersionsWrongRemote() {
	s.NotPanics(func() {
		s.packetConfig.URLType = "wrongURLType"
		out, _ := yaml.Marshal(s.packetConfig)
		packet, _ := New(out)
		err := packet.FetchVersions()
		s.Assert().Error(err)
		s.Assert().Contains(err.Error(), "failed to find wrongURLType remote")
	})
}

func (s *packetSuite) TestFetchVersionsWrongRemoteURL() {
	s.NotPanics(func() {
		s.packetConfig.URL = "https://example.com"
		out, _ := yaml.Marshal(s.packetConfig)
		packet, _ := New(out)
		err := packet.FetchVersions()
		s.Assert().Error(err)
		s.Assert().Contains(err.Error(), "failed to get user and repo from packet URL")
	})
}

func (s *packetSuite) TestFindVersion() {
	s.NotPanics(func() {
		err := s.packet.FetchVersions()
		s.Assert().NoError(err)
		latestVersion := s.packet.Versions[0].Version.String()
		version, ok := s.packet.FindVersion(latestVersion)
		s.Assert().Equal(true, ok)
		s.Assert().Equal(version.Version.String(), latestVersion)
	})
}

func (s *packetSuite) TestFindVersionFails() {
	s.NotPanics(func() {
		err := s.packet.FetchVersions()
		s.Assert().NoError(err)
		latestVersion := "WrongVersion"
		version, ok := s.packet.FindVersion(latestVersion)
		s.Assert().Equal(false, ok)
		s.Assert().Equal(version, Version{})
	})
}

func (s *packetSuite) TestLatestVersion() {
	s.NotPanics(func() {
		version, err := s.packet.LatestVersion()
		latestVersion := s.packet.Versions[0]
		s.Assert().NoError(err)
		s.Assert().Equal(version, latestVersion)
	})
}

func (s *packetSuite) TestEmptyVersionListLatestVersion() {
	s.Assert().NotPanics(func() {
		s.packetConfig.URL = "http://github.com/t/xxxg"
		out, _ := yaml.Marshal(s.packetConfig)
		packet, err := New(out)
		s.Assert().NoError(err)
		_, err = packet.LatestVersion()
		s.Assert().Error(err)
		s.Assert().Contains(err.Error(), "version list has 0 elements")
	})
}

func TestPackets(t *testing.T) {
	suite.Run(t, new(packetSuite))
}
