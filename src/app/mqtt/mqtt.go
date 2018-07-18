package mqtt

import (
	"fmt"
	"os/exec"
	"log"
)

func SendMQTT(device string, command string) {
	cmd := exec.Command(command)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed with error: ", err)
	}
	fmt.Println(string(out)) // just for testing
}
