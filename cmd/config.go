package cmd

import (
	"github.com/it-sova/bin-manager/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get user home dir, %w", err)
	}

	configDir := path.Join(home, ".config", "binm")
	installDir := path.Join(home, ".binm")
	stateLocation := path.Join(home, ".config", "binm", ".state.yaml")

	if err = helpers.CreateDirIfNotExists(configDir); err != nil {
		log.Error(err)
	}

	viper.SetDefault("GithubToken", "")
	viper.SetDefault("InstallDir", installDir)
	viper.SetDefault("StateLocation", stateLocation)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(configDir)
	err = viper.ReadInConfig()

	if err != nil {
		log.Infof("Failed to read config file, using defaults to operate: %v", err.Error())
	}

}
