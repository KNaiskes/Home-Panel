package main

import (
	"app/mqtt"
	"app/database"
)

func main() {
	go RunEvery(database.DBexists, 10)
	go GetDht22(10)
	mqtt.Dht22Mqtt()
}
