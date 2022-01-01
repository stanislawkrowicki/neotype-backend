package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

const (
	envKey             = "ENVIRONMENT"
	configFolder       = "config"
	defaultEnvironment = "local"
	envProd            = "production"
	envPort            = "PORT"
)

func Get(pkg, key string) (string, error) {
	configType, exists := os.LookupEnv(envKey)
	if !exists {
		configType = defaultEnvironment
	}

	mainPath, _ := os.Getwd()

	path := fmt.Sprintf("%s/%s/%s/config_%s.yaml", mainPath, configFolder, pkg, configType)

	if strings.Contains(mainPath, "tests") {
		path = fmt.Sprintf("%s/../../%s/%s/config_%s.yaml", mainPath, configFolder, pkg, configType)
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	data := make(map[string]string)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return "", err
	}

	if val, ok := data[key]; ok {
		return val, nil
	} else {
		return "", err
	}
}

func GetBaseURL(pkg string) (string, error) {
	baseURL, err := Get(pkg, "base")
	if err != nil {
		return "", err
	}

	if baseURL == "http://localhost" {
		port, err := Get(pkg, "port")
		if err != nil {
			return "", err
		}
		baseURL += ":" + port
	}

	return baseURL, nil
}

func GetPort(service string) (string, error) {
	environment, exists := os.LookupEnv(envKey)
	if !exists || environment != envProd {
		return Get(service, "port")
	}

	port, exists := os.LookupEnv(envPort)
	if !exists {
		return "", fmt.Errorf("PORT variable not set, even though in production mode")
	}

	return port, nil
}
