package cmd

import (
	"errors"
	"github.com/it-sova/bin-manager/helpers"
	"os"
	"path"
	"strings"

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
		// Check if $PATH includes installPath
		if !strings.Contains(os.Getenv("PATH"), installPath) {
			log.Warningf("Install path %v not found in $PATH!", installPath)
			log.Warningf("Consider extending $PATH: export PATH=\"$PATH:%v\"", installPath)
		}

		err := helpers.CreateDirIfNotExists(installPath)
		if err != nil {
			log.Fatalf("Failed to create %v, %v", installPath, err)
		}

		packets.Load()
		packet, err := packets.FindPacket(args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("Going to install packet ", packet.Name)
		err = packet.Install(installPath)
		if err != nil {
			log.Fatalf("Failed to install packet: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory, %v", err)
	}

	installCmd.Flags().StringVarP(
		&installPath,
		"path",
		"p",
		path.Join(homeDir, ".binm"),
		"Define path for packets installation",
	)
}
