package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/it-sova/bin-manager/packets"
	"github.com/it-sova/bin-manager/state"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateAll bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			updateAll = true
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		binState, err := state.Get()

		if err != nil {
			log.Fatalf(err.Error())
		}

		if updateAll {
			log.Info("Running complete packets update")
			for _, installedPacket := range binState.InstalledPackets {
				err := updatePacket(installedPacket, binState)
				if err != nil {
					log.Fatalf(err.Error())
				}
			}
		} else {
			packetName := args[0]
			log.Infof("Running %v update", packetName)

			if installedPacket, ok := binState.FindInstalledPacket(packetName, ""); ok {
				err := updatePacket(installedPacket, binState)
				if err != nil {
					log.Fatalf(err.Error())
				}
			} else {
				log.Fatalf("Failed to find %v in state,", packetName)
			}

		}
	},
}

func updatePacket(installedPacket state.InstalledPacket, binState state.State) error {
	log.Debugf("Packet -> %v", installedPacket)
	packet, err := packets.FindPacket(installedPacket.Name)

	if err != nil {
		log.Warningf("Packet %v found in state but not found in repo", installedPacket.Name)
	}

	latestPacketVersion, err := packet.LatestVersion()
	if err != nil {
		return err
	}

	installedPacketVersion, err := version.NewVersion(installedPacket.Version)
	if err != nil {
		return fmt.Errorf("failed to parse %v version from state", installedPacket.Name)
	}

	if installedPacketVersion.LessThan(latestPacketVersion.Version) {
		log.Infof("%v version %v less then available %v, updating...", installedPacket.Name, installedPacket.Version, latestPacketVersion.Version.String())

		installPath := path.Dir(installedPacket.Path)
		err = packet.Install(installPath, latestPacketVersion.Version.String())

		if err != nil {
			return fmt.Errorf("failed to reinstall packet, %v", err.Error())
		}

		err = binState.Remove(installedPacket.Name)
		if err != nil {
			return fmt.Errorf("failed to remove updated packet from state, %v", err.Error())
		}

		err = binState.Append(state.InstalledPacket{
			Name:      packet.Name,
			Version:   latestPacketVersion.Version.String(),
			Path:      installedPacket.Path,
			Installed: time.Now().String(),
		})

		if err != nil {
			return fmt.Errorf("failed to save state, %v", err)
		}

		log.Infof("Finished updating %v", installedPacket.Name)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
