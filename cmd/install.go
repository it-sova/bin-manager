package cmd

import (
	"errors"
	"fmt"
	"github.com/it-sova/bin-manager/helpers"
	"github.com/it-sova/bin-manager/state"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
	"time"

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

		if len(installPath) == 0 {
			log.Debug("Install path not defined via CLI, use path from config or default")
			installPath = viper.Get("InstallDir").(string)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			log.Error("Failed to get user home dir, %w", err)
		}

		// Replace ~/... -> /home/user/...
		installPath = strings.Replace(installPath, "~", home, 1)

		// Check if $PATH includes installPath
		if !strings.Contains(os.Getenv("PATH"), installPath) {
			log.Warningf("Install path %v not found in $PATH!", installPath)
			log.Warningf("Consider extending $PATH: export PATH=\"$PATH:%v\"", installPath)
		}

		err = helpers.CreateDirIfNotExists(installPath)
		if err != nil {
			log.Fatalf("Failed to create %v, %v", installPath, err)
		}

		packets.Load()

		packet, err := packets.FindPacket(args[0])
		if err != nil {
			log.Fatal(err)
		}

		binState, err := state.Get()
		if err != nil {
			fmt.Errorf("failed to get state, %v", err)
		}

		if installedPacket, ok := binState.FindPacket(packet.Name); ok {
			log.Infof("Packet %v %v already installed", installedPacket.Name, installedPacket.Version)
			os.Exit(0)
		}

		log.Debug("Going to install packet ", packet.Name)
		err = packet.Install(installPath, "")
		if err != nil {
			log.Fatalf("Failed to install packet: %v", err)
		}

		err = binState.Append(state.InstalledPacket{
			Name:      packet.Name,
			Version:   packet.Versions[0].Version.String(),
			Path:      path.Join(installPath, packet.Name),
			Installed: time.Now().String(),
		})

		if err != nil {
			log.Errorf("failed to save state, %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(
		&installPath,
		"path",
		"p",
		"",
		"Define path for packets installation",
	)
}
