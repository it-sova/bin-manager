package packets

import (
	"fmt"

	"github.com/it-sova/bin-manager/repo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var packets []packet

func Load() {
	repos := repo.RepoList()
	for _, repo := range repos {
		log.Info(fmt.Sprintf("Loaded repo %v (%v)", repo.GetName(), repo.GetPath()))

		for _, packetName := range repo.ScanPackets() {
			log.Debug("Found packet ", packetName)

			packetConfig, err := repo.GetPacketConfig(packetName)
			if err != nil {
				log.Error("Failed to read packet config for ", packetName)
			}

			packet, err := NewPacket(packetConfig)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to unmarshal packet config for %v: %v", packetName, err.Error()))
			}

			log.Debug(fmt.Sprintf("Packet loaded: %s", packet.Name))
			packets = append(packets, packet)
		}
	}
}

func NewPacket(config []byte) (packet, error) {
	packet := packet{}
	err := yaml.Unmarshal(config, &packet)
	if err != nil {
		return packet, err
	}

	return packet, nil
}

func ListAll() {
	log.Printf("%+v", packets)
}

func FindPacket(name string) (packet, error) {
	if len(packets) == 0 {
		Load()
	}

	for _, packet := range packets {
		if packet.Name == name {
			return packet, nil
		}
	}

	return packet{}, fmt.Errorf("Unable to find packet %v", name)

}
