package cmd

import (
	"errors"
	"fmt"
	"github.com/it-sova/bin-manager/state"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/it-sova/bin-manager/packets"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Single packet name required")
		}
		return nil
	},
	Aliases: []string{
		"uninstall",
		"delete",
	},
	Run: func(cmd *cobra.Command, args []string) {
		packets.Load()

		packet, err := packets.FindPacket(args[0])
		if err != nil {
			log.Fatal(err)
		}

		binState, err := state.Get()
		if err != nil {
			fmt.Errorf("failed to get state, %v", err)
		}

		installedPacket, ok := binState.FindInstalledPacket(packet.Name, "")
		if !ok {
			log.Fatalf("Packet %v not found in state", packet.Name)
		}

		log.Infof("Uninstalling %v %v (%v)", installedPacket.Name, installedPacket.Version, installedPacket.Path)

		err = os.Remove(installedPacket.Path)
		if err != nil {
			log.Warningf("Failed to remove installed packet (%v), removing from state", err)
		}

		err = binState.Remove(packet.Name)
		if err != nil {
			log.Fatalf("Failed to update state, %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
