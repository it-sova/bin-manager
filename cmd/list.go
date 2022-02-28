package cmd

import (
	"github.com/it-sova/bin-manager/packets"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List command",
	Run: func(cmd *cobra.Command, args []string) {
		packets.Load()
		packets := packets.GetAll()
		for _, packet := range packets {
			packet.ListVersions()
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
