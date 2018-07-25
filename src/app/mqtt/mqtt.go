package mqtt

import (
	"os/exec"
)

func ChangeState(command string, topic string) {
	// function for every two state (on/off) electrical devices
	switch command {
		case "on":
			cmd := exec.Command("mosquitto_pub", "-m", "on", "-t", topic)
			cmd.Output()
		case "off":
			cmd := exec.Command("mosquitto_pub", "-m", "off", "-t", topic)
			cmd.Output()
	}
}

func ChangeColor(color string, topic string) {
	// function for led strips and bubles rgb bulbs
	cmd := exec.Command("mosquitto_pub", "-m", color, "-t", topic)
	cmd.Output()
}


