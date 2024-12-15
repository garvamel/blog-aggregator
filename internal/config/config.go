package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFilename = "/.gatorconfig.json"
const projectPath = "/Workspace/blog-aggregator"

func getConfigFilePath() (string, error) {
	// Construct the file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// log.Fatal(err)
		return "", err
	}
	return homeDir + projectPath + configFilename, nil
}

func Read() Config {

	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	configJson := Config{}

	err = json.Unmarshal(data, &configJson)
	if err != nil {
		log.Fatal(err)
	}

	// Print the contents
	return configJson
}

func (c *Config) SetUser(user string) {

	c.CurrentUserName = user

	json, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(filePath, json, 0777)
}
