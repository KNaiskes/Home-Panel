package mqtt

import (
	"os/exec"
)

func SendMQTT(device string) {

	switch device {
		case "light1":
			cmd := exec.Command("mosquitto_pub", "-m", "light1", "-t", "test")
			cmd.Output()
		case "light2":
			cmd := exec.Command("mosquitto_pub", "-m", "light2", "-t", "test")
			cmd.Output()

		}
	}

