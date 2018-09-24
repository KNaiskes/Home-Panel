package main

import (
	"github.com/KNaiskes/Home-Panel/mqtt"
	"github.com/KNaiskes/Home-Panel/database"
)

func main() {
	go RunEvery(database.DBexists, 10)
	go GetDht22(10)
	mqtt.Dht22Mqtt()
}
