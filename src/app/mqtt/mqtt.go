package mqtt

import (
	"os/exec"
)

func ChangeState(command string, topic string) {

	switch command {
		case "on":
			cmd := exec.Command("mosquitto_pub", "-m", "on", "-t", topic)
			cmd.Output()
		case "off":
			cmd := exec.Command("mosquitto_pub", "-m", "off", "-t", topic)
			cmd.Output()
	}

}
