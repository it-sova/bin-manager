package packets

import (
	"fmt"

	"github.com/it-sova/bin-manager/repo"
	log "github.com/sirupsen/logrus"
)

var packets []Packet

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

func ListAll() {
	log.Printf("%+v", packets)
}

func FindPacket(name string) (Packet, error) {
	if len(packets) == 0 {
		Load()
	}

	for _, packet := range packets {
		if packet.Name == name {
			return packet, nil
		}
	}

	return Packet{}, fmt.Errorf("Unable to find packet %v", name)

}
