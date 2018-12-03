package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	username	string
	password	string
	endpoint	string
	database	string
}

func parseConfigFile(configFile string) Configuration {
	file, _ := os.Open(configFile)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return configuration
}
