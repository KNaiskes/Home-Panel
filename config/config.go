package config

import(
	"encoding/json"
	"os"
	"log"
)

type Config struct {
	ServerIP string `json:"server_ip"`
	ClientID string `json:"client_id"`
}

func Getconfig() Config {
	file, err := os.Open("config.json")

	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	config := Config{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
