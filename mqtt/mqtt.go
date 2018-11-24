package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/KNaiskes/Home-Panel/database"
	"github.com/KNaiskes/Home-Panel/config"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"strings"
)

var myServer = config.Getconfig().ServerIP
var clientId = config.Getconfig().ClientID
var MqttUsername = config.Getconfig().MqttUsername
var MqttPassword = config.Getconfig().MqttPassword

func ChangeState(command string, topic string) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(myServer)
	opts.SetClientID(clientId)
	opts.SetUsername(MqttUsername)
	opts.SetPassword(MqttPassword)
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
	opts.SetUsername(MqttUsername)
	opts.SetPassword(MqttPassword)
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
	//temperature, humidity := split[0], split[1]

	if metrics != "getTempHum" {
		temperature, _ := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
		humidity, _ := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)

		if temperature > 0  && humidity > 0 {
			database.AddTempHum(temperature, humidity)
		}
	}
}

func Dht22Mqtt() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	topic := "SensorDht22"
	qos := 0
	clientid := "clientid"

	connOpts := MQTT.NewClientOptions().AddBroker(myServer).SetClientID(clientid).SetUsername(MqttUsername).SetPassword(MqttPassword).SetCleanSession(true)

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
	opts.SetUsername(MqttUsername)
	opts.SetPassword(MqttPassword)
	opts.SetCleanSession(true)

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := c.Publish(topic, 0, false, command)
	token.Wait()
	c.Disconnect(250)

}
