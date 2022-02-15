package cmd

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/it-sova/bin-manager/packets"
	"github.com/spf13/cobra"
)

var installPath string

// installCmd represents the list command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Single packet name required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		packets.Load()
		packet, err := packets.FindPacket(args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("Going to install packet ", packet.Name)
		packet.Install()

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installPath, "path", "p", "/opt/binm/", "Define path for packets installation")
	//installCmd.Flags().StringP("path", "p", "/opt/binm/", "Define path for packets installation")
}
