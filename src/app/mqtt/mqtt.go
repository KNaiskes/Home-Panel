package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func ChangeState(command string, topic string) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://192.168.1.30:1883")
	opts.SetClientID("mqttServer")
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
	opts.AddBroker("tcp://192.168.1.30:1883")
	opts.SetClientID("mqttServer")
	opts.SetCleanSession(true)

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := c.Publish(topic, 0, false, color)
	token.Wait()
	c.Disconnect(250)
}
