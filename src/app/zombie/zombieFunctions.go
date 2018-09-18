package main

import (
	"app/mqtt"
	"time"
)

func GetDht22(minutes time.Duration) {
	topic := "SensorDht22"
	command := "getTempHum"

	for {
		<-time.After(minutes * time.Second) //seconds instead of minutes just for testing
		go mqtt.SendMsg(topic, command)
	}
}

func RunEvery(f func(), minutes time.Duration) {
	for {
		<-time.After(minutes * time.Second)
		go f()
	}
}
