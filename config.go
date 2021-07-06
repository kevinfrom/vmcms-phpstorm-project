package main

import (
	"fmt"
)

type config struct {
	parsedConfig      map[string]string
	configPath        string
	configExamplePath string
}

/*
** Get config
 */
func getConfig(configPath string, configExamplePath string) config {
	config := config{
		parsedConfig:      map[string]string{},
		configPath:        configPath,
		configExamplePath: configExamplePath,
	}
	config.parseConfig()

	return config
}

/*
** Parse config, complete drive_path and phpstorm_path
 */
func (config config) parseConfig() {
	config.parseConfigFile()

	exitIfFileDoesNotExist(config.parsedConfig["drive_path"])
	exitIfFileDoesNotExist(config.parsedConfig["projects_path"])
	exitIfFileDoesNotExist(config.parsedConfig["phpstorm_path"])

	config.parsedConfig["phpstorm_path"] = findPhpstormExecutable(config.parsedConfig["phpstorm_path"])
}

/*
** Checks if config.json file has all the necessary values
 */
func (config config) parseConfigFile() {
	exitIfFileDoesNotExist(config.configPath)
	exitIfFileDoesNotExist(config.configExamplePath)

	configJson := readJsonFile(config.configPath)
	configExampleJson := readJsonFile(config.configExamplePath)

	// Check length of maps
	errorMessage := fmt.Sprintf("%s does not match pattern of %s", config.configPath, config.configExamplePath)
	if len(configJson) != len(configExampleJson) {
		exitWithErrorMessage(errorMessage)
	}

	// Check if config keys are the same and values are not empty
	for exampleKey := range configExampleJson {
		isError := false

		value, keyExists := configJson[exampleKey]

		if keyExists == false {
			isError = true
		} else if len(value) <= 0 {
			isError = true
		}

		if isError {
			exitWithErrorMessage(errorMessage)
		}
	}

	// Set new values on config struct
	for key, value := range configJson {
		config.parsedConfig[key] = value
	}
}
