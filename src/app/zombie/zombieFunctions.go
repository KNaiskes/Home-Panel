package main

import (
	"app/mqtt"
	"time"
)

func GetDht22(minutes time.Duration) {
	topic := "SensorDht22"
	getTemperature := "getTemp"
	getHumidity := "getHum"

	for {
		<-time.After(minutes * time.Second) //seconds instead of minutes just for testing
		go mqtt.SendMsg(topic, getTemperature)
		time.Sleep(2 * time.Second)
		go mqtt.SendMsg(topic, getHumidity)
	}
}
