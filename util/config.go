package util

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	QueueName string `json:"queue_name"`
}

func ReadConfig() *Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	var config Config
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatal("Parsing config file", err)
	}
	return &config
}
