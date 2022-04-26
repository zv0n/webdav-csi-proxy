package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	SocketPath string
}

func ParseConfigFile(configPath string) (*Configuration, error) {
	// prepare a default configuration
	configuration := Configuration{
		SocketPath: "/tmp/webdav-proxy.socket",
	}

	file, err := os.Open(configPath)
	if err != nil {
		return &configuration, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return &configuration, err
	}

	return &configuration, nil
}
