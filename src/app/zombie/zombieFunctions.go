package main

import (
	"app/mqtt"
	"time"
)

func GetDht22(minutes time.Duration , topic string, command string) {
	for {
		<-time.After(minutes * time.Second)
		go mqtt.SendMsg(topic, command)
	}
}
