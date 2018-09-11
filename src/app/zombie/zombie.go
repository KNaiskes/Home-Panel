package main

import "app/mqtt"

func main() {
	go GetDht22(10, "SensorDht22", "getTemp")
	//go GetDht22(12, "SensorDht22", "getHum")
	mqtt.Dht22Mqtt()
}
