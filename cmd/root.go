package cmd

import (
	"os"
	"path"

	"github.com/it-sova/bin-manager/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bin-manager",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "Log level")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get user home dir, %w", err)
	}

	configDir := path.Join(home, ".config", "binm")

	if err = helpers.CreateDirIfNotExists(configDir); err != nil {
		log.Error(err)
	}

	// TODO: Get GitHub token from env or config
	viper.SetDefault("InstallDir", "/opt/binm/")

	viper.SetConfigName("binm")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(configDir)
	err = viper.ReadInConfig()

	if err != nil {
		log.Infof("Failed to read config file, using defaults to operate: %v", err.Error())
	}

	log.Infof("Installation directory: %v", viper.Get("InstallDir"))
}
