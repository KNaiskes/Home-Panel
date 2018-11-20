package config

import(
	"encoding/json"
	"os"
	"log"
)

const filePath = "src/github.com/KNaiskes/Home-Panel/config.json"

type Config struct {
	ServerIP string `json:"server_ip"`
	ClientID string `json:"client_id"`
}

func Getconfig() Config {
	file, err := os.Open(filePath)

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
