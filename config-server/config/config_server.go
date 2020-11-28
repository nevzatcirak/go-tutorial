package config

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type Property map[string]interface{}

type AppProperties struct {
	property *Property
}

var instance *AppProperties

func GetProperty(key string) string {
	return fmt.Sprintf("%v", (*instance.property)[key])
}

func GetAllProperties() *Property {
	return instance.property
}

func LoadProperties() {
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

	instance = &AppProperties{
		property: mergeProperties(&config, appName),
	}
}

func mergeProperties(configResource *cloudconfigclient.Source, appName string) *Property {
	property := Property{}
	sortedProperties := getSortedProperties(configResource, appName)
	for i := range *sortedProperties {
		for j, val := range (*sortedProperties)[i].Source {
			property[j] = val
		}
	}
	fillExactKeyValue(&property)
	return &property
}

func fillExactKeyValue(property *Property) {
	for key, val := range *property {
		putValueIntoPlace(key, val, property)
	}
}

func putValueIntoPlace(key string, value interface{}, property *Property) {
	val := fmt.Sprintf("%v", value)
	if strings.Index(val, "${") >= 0 && strings.Index(val, "}") >= 0 {
		bareKey := val[strings.Index(val, "${")+2 : strings.Index(val, "}")]
		keyWithSign := val[strings.Index(val, "${") : strings.Index(val, "}")+1]
		exactValue := fmt.Sprintf("%v", (*property)[bareKey])
		putValueIntoPlace(bareKey, exactValue, property)
		exactValue = fmt.Sprintf("%v", (*property)[bareKey])
		(*property)[key] = strings.Replace(val, keyWithSign, exactValue, 1)
		putValueIntoPlace(key, (*property)[key], property)
	}
}

func getSortedProperties(configResource *cloudconfigclient.Source, appName string) *[]cloudconfigclient.PropertySource {
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
			if (strings.Contains(propertyName, "application.yml") ||
				strings.Contains(propertyName, "application.properties")) &&
				strings.EqualFold("application", propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			} else if profile != "default" &&
				(strings.Contains(propertyName, "application-"+profile+".yml") ||
					strings.Contains(propertyName, "application-"+profile+".properties")) &&
				strings.EqualFold("application-"+profile, propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			} else if (strings.Contains(propertyName, appName+".yml") ||
				strings.Contains(propertyName, appName+".properties")) &&
				strings.EqualFold(appName, propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			} else if profile != "default" &&
				(strings.Contains(propertyName, appName+"-"+profile+".yml") ||
					strings.Contains(propertyName, appName+"-"+profile+".properties")) &&
				strings.EqualFold(appName+"-"+profile, propNames[j]) {
				targetPropertySource = append(targetPropertySource, property)
				break
			}
		}
	}
	return &targetPropertySource
}
