package packets

import (
	"fmt"

	"github.com/it-sova/bin-manager/repo"
	log "github.com/sirupsen/logrus"
)

var packets []Packet

//TODO: Add ability to load single packet
// Load loads all packets from all repos into packets slice
func Load() {
	repos := repo.List()
	for _, repo := range repos {
		log.Info(fmt.Sprintf("Loaded repo %v (%v)", repo.GetName(), repo.GetPath()))

		for _, packetName := range repo.ScanPackets() {
			log.Debug("Found packet ", packetName)

			packetConfig, err := repo.GetPacketConfig(packetName)
			if err != nil {
				log.Error("Failed to read packet config for ", packetName)
			}

			packet, err := New(packetConfig)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to unmarshal packet config for %v: %v", packetName, err.Error()))
			}

			log.Debug(fmt.Sprintf("Packet loaded: %s", packet.Name))
			packets = append(packets, packet)
		}
	}
}

func GetAll() []Packet {
	return packets
}

// ListAll lists all loaded packets
func ListAll() {
	for _, packet := range packets {
		log.Infof("%v - %v - (%v)", packet.Name, packet.Description, packet.URL)
	}
}

// FindPacket finds packet in loaded list by its name
func FindPacket(name string) (Packet, error) {
	if len(packets) == 0 {
		Load()
	}

	for _, packet := range packets {
		if packet.Name == name {
			return packet, nil
		}
	}

	return Packet{}, fmt.Errorf("unable to find packet %v", name)
}

// ListVersions parses remote to get available packet versions
func (p *Packet) ListVersions() {
	if err := p.FetchVersions(); err == nil {
		log.Infof("= %v:", p.Name)

		for _, v := range p.Versions {
			log.Infof("    %v", v.Version)
		}
	} else {
		log.Errorf("Failed to fetch packet versions, %v", err.Error())
	}
}
