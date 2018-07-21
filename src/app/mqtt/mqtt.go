package mqtt

import (
	"os/exec"
)

func SendMQTT(command string) {

	switch command {
		case "on":
			cmd := exec.Command("mosquitto_pub", "-m", "on", "-t", "ledStrip")
			cmd.Output()
		case "off":
			cmd := exec.Command("mosquitto_pub", "-m", "off", "-t", "ledStrip")
			cmd.Output()
	}

}
