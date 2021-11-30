package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

const (
	envKey = "ENVIRONMENT"
	configFolder = "config"
	defaultEnvironment = "local"
)
func Get(pkg, key string) (string, error) {
	configType, exists := os.LookupEnv(envKey)
	if !exists {
		configType = defaultEnvironment
	}

	path := fmt.Sprintf("%s/%s/config_%s.yaml", configFolder, pkg, configType)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	data := make(map[string]string)
	err = yaml.Unmarshal(file, &data); if err != nil {
		return "", err
	}

	if val, ok := data[key]; ok {
		return val, nil
	} else {
		return "", err
	}
}
