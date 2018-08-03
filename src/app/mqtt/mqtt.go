package mqtt

import (
	"fmt"
	"os/exec"
	"os"
	"time"
	"log"
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
	// function for led strips and rgb bulbs
	cmd := exec.Command("mosquitto_pub", "-m", color, "-t", topic)
	cmd.Output()
}

func RunEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func Dht22Sensor(t time.Time) {
	tempFile := "temperatureHUm.txt"
	topic := "tempHum"
	cmd := exec.Command("mosquitto_sub", "-t", topic)
	outfile, err := os.Create(tempFile)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile
	fmt.Println("here!!!")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(20 * time.Second)
	os.Remove(tempFile)

}



