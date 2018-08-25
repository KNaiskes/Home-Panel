package mqtt


import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"os/exec"
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

	//token := c.Publish("ledStrip", 0, false, "on")
	//token.Wait()
	//c.Disconnect(250)
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


//import (
//	"os/exec"
//)
//
//func ChangeState(command string, topic string) {
//	// function for every two state (on/off) electrical devices
//	switch command {
//		case "on":
//			cmd := exec.Command("mosquitto_pub", "-m", "on", "-t", topic)
//			cmd.Output()
//		case "off":
//			cmd := exec.Command("mosquitto_pub", "-m", "off", "-t", topic)
//			cmd.Output()
//	}
//}
//
//func ChangeColor(color string, topic string) {
//	// function for led strips and rgb bulbs
//	cmd := exec.Command("mosquitto_pub", "-m", color, "-t", topic)
//	cmd.Output()
//}
