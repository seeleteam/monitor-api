package cmd

import (
	"encoding/json"
	"io/ioutil"
)

// GetConfigFromFile unmarshals the config from the given file
func GetConfigFromFile(filepath string) (map[string]string, error) {
	var config map[string]string
	buff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(buff, &config)
	return config, err
}
