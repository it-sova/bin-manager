package state

import (
	"fmt"
	"io"
	"os"

	"github.com/it-sova/bin-manager/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// InstalledPacket represents installed packet as part of state
type InstalledPacket struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	Path      string `yaml:"path"`
	Installed string `yaml:"installed"`
}

// State represent state ( list of installed packets )
type State struct {
	InstalledPackets []InstalledPacket `yaml:"installed_packets"`
	location         string
}

// Get return state. If state file was not found - it will be created with empty state
func Get() (State, error) {
	stateLocation, ok := viper.Get("StateLocation").(string)
	if !ok || stateLocation == "" {
		return State{}, fmt.Errorf("failed to get state location")
	}

	if _, err := os.Stat(stateLocation); err != nil {
		log.Infof("Creating empty state at %v", stateLocation)

		err := createEmptyState()
		if err != nil {
			return State{}, fmt.Errorf("failed to create new state, %v", err)
		}
	}

	stateFile, err := os.Open(stateLocation)
	if err != nil {
		return State{}, fmt.Errorf("failed to open state file at %v, %v", stateFile, err)
	}
	defer stateFile.Close()
	data, err := io.ReadAll(stateFile)

	if err != nil {
		return State{}, fmt.Errorf("failed to read state file at %v, %v", stateFile, err)
	}

	state := State{
		location: stateLocation,
	}
	err = yaml.Unmarshal(data, &state)

	if err != nil {
		return State{}, fmt.Errorf("failed to unmarshal state file, %v", err)
	}

	return state, nil
}

// createEmptyState creates empty state file
func createEmptyState() error {
	stateLocation, ok := viper.Get("StateLocation").(string)

	if !ok || stateLocation == "" {
		return fmt.Errorf("failed to get state location")
	}

	out, err := yaml.Marshal(State{
		InstalledPackets: []InstalledPacket{},
	})

	if err != nil {
		return fmt.Errorf("failed to marshal empty state")
	}

	err = os.WriteFile(stateLocation, out, helpers.FileChmod)

	if err != nil {
		return fmt.Errorf("failed to save empty state at %v", stateLocation)
	}

	return nil
}

// Append appends new InstalledPacket into state
func (s *State) Append(packet InstalledPacket) error {
	// Insert element to beginning of slice
	s.InstalledPackets = append([]InstalledPacket{packet}, s.InstalledPackets...)
	err := s.Save()

	return err
}

// Remove removes installed packet from state by packet name
func (s *State) Remove(packetName string) error {
	for index, packet := range s.InstalledPackets {
		if packet.Name == packetName {
			s.InstalledPackets = append(s.InstalledPackets[:index], s.InstalledPackets[index+1:]...)
		}
	}

	log.Debugf("State - %#v", s.InstalledPackets)
	err := s.Save()

	return err
}

// Save saves state to file
func (s *State) Save() error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state, %v", err)
	}

	err = os.WriteFile(s.location, data, helpers.FileChmod)
	if err != nil {
		return fmt.Errorf("failed to write state file, %v", err)
	}

	return nil
}

// FindInstalledPacket searches for installed packet in state. It's possible to search by name or name and version
func (s *State) FindInstalledPacket(name, version string) (InstalledPacket, bool) {
	for _, packet := range s.InstalledPackets {
		if version != "" {
			if packet.Name == name && packet.Version == version {
				return packet, true
			}
		} else {
			if packet.Name == name {
				return packet, true
			}
		}
	}

	return InstalledPacket{}, false
}
