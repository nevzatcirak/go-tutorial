package config

import (
	"github.com/Piszmog/cloudconfigclient"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	Strings "strings"
	"time"
)

type Property map[string]interface{}

func GetProperties() *[]cloudconfigclient.PropertySource {
	configServerUrl := os.Getenv("CONFIG_SERVER_URL")
	appName := os.Getenv("APPLICATION_NAME")
	profile := os.Getenv("PROFILE")
	if configServerUrl == "" {
		log.Fatal("CONFIG_SERVER_URL is not defined as environment variable.")
	} else if appName == "" {
		log.Fatal("APPLICATION_NAME is not defined as environment variable.")
	} else if profile == "" {
		log.Warn("PROFILE is not defined as environment variable. PROFILE will be set as 'default'.")
		profile = "default"
	}
	configClient, err := cloudconfigclient.NewLocalClient(&http.Client{}, []string{configServerUrl})

	if err != nil {
		panic(err)
	}

	// Retrieves the configurations from the Config Server based on the application name and active profiles
	config, err := configClient.GetConfiguration(appName, []string{profile})
	for err != nil {
		log.Error(err)
		time.Sleep(2 * time.Second)
		config, err = configClient.GetConfiguration(appName, []string{profile})
	}

	log.Info(mergeProperties(&config, appName, profile))
	return &config.PropertySources
}

func mergeProperties(configResource *cloudconfigclient.Source, appName string, profile string) *Property {
	property := Property{}
	sortedProperties := getSortedProperties(*configResource, appName)
	for i := range *sortedProperties {
		for j, val := range (*sortedProperties)[i].Source {
			property[j] = val
		}
	}
	return &property
}

func getSortedProperties(configResource cloudconfigclient.Source, appName string) *[]cloudconfigclient.PropertySource {
	currentPropertySource := configResource.PropertySources
	var targetPropertySource []cloudconfigclient.PropertySource
	profile := configResource.Profiles[0]
	var propNames []string
	if profile == "default" {
		propNames = []string{"application", appName}
	} else {
		propNames = []string{"application", "application-" + profile, appName, appName + "-" + profile}
	}

	for j := 0; j < len(propNames); j++ {
		for _, property := range currentPropertySource {
			propertyName := property.Name
			if (Strings.Contains(propertyName, "application.yml") ||
				Strings.Contains(propertyName, "application.properties")) &&
				Strings.EqualFold("application", propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			} else if profile != "default" &&

				(Strings.Contains(propertyName, "application-"+profile+".yml") ||
					Strings.Contains(propertyName, "application-"+profile+".properties")) &&
				Strings.EqualFold("application-"+profile, propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			} else if (Strings.Contains(propertyName, appName+".yml") ||
				Strings.Contains(propertyName, appName+".properties")) &&
				Strings.EqualFold(appName, propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			} else if profile != "default" &&
				(Strings.Contains(propertyName, appName+"-"+profile+".yml") ||
					Strings.Contains(propertyName, appName+"-"+profile+".properties")) &&
				Strings.EqualFold(appName+"-"+profile, propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			}
		}
	}
	return &targetPropertySource
}
