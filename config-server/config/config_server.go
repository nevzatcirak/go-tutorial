package config

import (
	"github.com/Piszmog/cloudconfigclient"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func GetProperties() *[]cloudconfigclient.PropertySource {
	configServerUrl := os.Getenv("CONFIG_SERVER_URL")
	appName := os.Getenv("APPLICATION_NAME")
	profile := os.Getenv("PROFILE")
	if configServerUrl == "" {
		log.Fatal("CONFIG_SERVER_URL is not defined as environment variable.")
	} else if appName == "" {
		log.Fatal("APPLICATION_NAME is not defined as environment variable.")
	} else if profile == "" {
		log.Fatal("PROFILE is not defined as environment variable.")
	}
	configClient, err := cloudconfigclient.NewLocalClient(&http.Client{}, []string{configServerUrl})

	if err != nil {
		panic(err)
	}

	// Retrieves the configurations from the Config Server based on the application name and active profiles
	config, err := configClient.GetConfiguration(appName, []string{profile})
	if err != nil {
		panic(err)
	}

	return &config.PropertySources
}
