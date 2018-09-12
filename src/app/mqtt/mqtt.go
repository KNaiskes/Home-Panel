package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"app/database"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"strings"
)

const myServer = "tcp://192.168.1.30:1883"
const clientId = "mqttServer"

func ChangeState(command string, topic string) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(myServer)
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	switch command {
	case "on":
		token := c.Publish(topic, 0, false, "on")
		token.Wait()
		c.Disconnect(250)
	case "off":
		token := c.Publish(topic, 0, false, "off")
		token.Wait()
		c.Disconnect(250)
	}
}

func ChangeColor(color string, topic string) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(myServer)
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := c.Publish(topic, 0, false, color)
	token.Wait()
	c.Disconnect(250)
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	var temporary  = message.Payload()

	var metrics string = string(temporary)

	split := strings.Split(metrics, ",")
	temperature, humidity := split[0], split[1]

	temperature, _ := strconv.ParseFloat(strings.TrimSpace(metrics), 64)
	humidity, _ := strconv.ParseFloat(strings.TrimSpace(metrics), 64)

	if metric > 0 {
		database.AddTempHum(temperature, humidity)

		fmt.Println("Metric: ", metric)
	}
}

func Dht22Mqtt() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	topic := "SensorDht22"
	qos := 0
	clientid := "clientid"

	connOpts := MQTT.NewClientOptions().AddBroker(myServer).SetClientID(clientid).SetCleanSession(true)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, byte(qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	<-c
}

func SendMsg(topic string, command string) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(myServer)
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := c.Publish(topic, 0, false, command)
	token.Wait()
	c.Disconnect(250)

}
