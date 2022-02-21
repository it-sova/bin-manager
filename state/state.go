package state

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type InstalledPacket struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	Path      string `yaml:"path"`
	Installed string `yaml:"installed"`
}

type State struct {
	InstalledPackets []InstalledPacket `yaml:"installed_packets"`
	location         string
}

func Get() (State, error) {
	stateLocation := viper.Get("StateLocation").(string)
	if stateLocation == "" {
		return State{}, fmt.Errorf("failed to get state location")
	}

	if _, err := os.Stat(stateLocation); err != nil {
		log.Debugf("State doesnt exists at %v, let's create empty one", stateLocation)
		err := CreateEmptyState()
		if err != nil {
			return State{}, fmt.Errorf("failed to create new state, %v", err)
		}
	}

	stateFile, err := os.Open(stateLocation)
	if err != nil {
		return State{}, fmt.Errorf("failed to open state file at %v, %v", stateFile, err)
	}
	defer stateFile.Close()
	data, err := ioutil.ReadAll(stateFile)
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

func CreateEmptyState() error {
	stateLocation := viper.Get("StateLocation").(string)
	if stateLocation == "" {
		return fmt.Errorf("failed to get state location")
	}

	out, err := yaml.Marshal(State{
		InstalledPackets: []InstalledPacket{},
	})

	if err != nil {
		return fmt.Errorf("failed to marshal empty state")
	}

	err = ioutil.WriteFile(stateLocation, out, 0644)

	if err != nil {
		return fmt.Errorf("failed to save empty state at %v", stateLocation)
	}

	return nil
}

func (s *State) Append(packet InstalledPacket) error {
	s.InstalledPackets = append(s.InstalledPackets, packet)
	err := s.Save()
	return err
}

func (s *State) Save() error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state, %v", err)
	}

	err = ioutil.WriteFile(s.location, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file, %v", err)
	}

	return nil

}

func (s *State) FindPacket(name string) (InstalledPacket, bool) {
	for _, packet := range s.InstalledPackets {
		if packet.Name == name {
			return packet, true
		}
	}

	return InstalledPacket{}, false

}
