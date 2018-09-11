package main

import "app/mqtt"

func main() {
	go GetDht22(10)
	mqtt.Dht22Mqtt()
}
