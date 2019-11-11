package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/MajorBreakfast/eventstore-user-operator/pkg/controller/eventstoreuser/eventstore"
)

type mainConfig struct {
	EventStores []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"eventStores"`
}

func loadEventStoreURLFromMainConfig(eventStore string) (string, error) {
	basePath := os.Getenv("EVENTSTORE_USER_OPERATOR_BASE_PATH")
	data, err := ioutil.ReadFile(path.Join(basePath, "config/config.yaml"))
	if err != nil {
		return "", err
	}

	loadedMainConfig := &mainConfig{}
	if err := yaml.Unmarshal([]byte(data), &loadedMainConfig); err != nil {
		return "", err
	}

	for _, item := range loadedMainConfig.EventStores {
		if item.Name == eventStore {
			return item.URL, nil
		}
	}

	return "", errors.New("Event Store config entry not found")
}

func loadEventStoreAdminUsername(eventStore string) (string, error) {
	basePath := os.Getenv("EVENTSTORE_USER_OPERATOR_BASE_PATH")
	data, err := ioutil.ReadFile(path.Join(basePath, "eventstore-credentials", eventStore, "username"))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func loadEventStoreAdminPassword(eventStore string) (string, error) {
	basePath := os.Getenv("EVENTSTORE_USER_OPERATOR_BASE_PATH")
	data, err := ioutil.ReadFile(path.Join(basePath, "eventstore-credentials", eventStore, "password"))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// LoadEventStoreConnection loads the config for an Event Store
func LoadEventStoreConnectionOptions(eventStore string) (eventstore.ConnectionOptions, error) {
	url, err := loadEventStoreURLFromMainConfig(eventStore)
	if err != nil {
		return eventstore.ConnectionOptions{}, err
	}

	adminUsername, err := loadEventStoreAdminUsername(eventStore)
	if err != nil {
		return eventstore.ConnectionOptions{}, err
	}

	adminPassword, err := loadEventStoreAdminPassword(eventStore)
	if err != nil {
		return eventstore.ConnectionOptions{}, err
	}

	return eventstore.ConnectionOptions{
		URL:           url,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
	}, nil
}
