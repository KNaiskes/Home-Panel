package mqtt

import (
	"fmt"
	"os/exec"
	"os"
	"time"
	"log"
	"bufio"
	"strconv"
	"strings"
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

func MakeFloats(textFile string) {
	floats := []float64 {}
	file, _ := os.Open(textFile)
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		trimmed := strings.TrimSpace(fileScanner.Text())
		f, er := strconv.ParseFloat(trimmed, 64)
		if er == nil {
			floats = append(floats, f)
		}
	}
		//TODO return slice or saving it to db instead of printing it
		fmt.Println("Floats: ", floats)
}

func Dht22Sensor(t time.Time) {
	tempFile := "temperatureHum.txt"
	topic := "tempHum"
	MakeFloats(tempFile)
	cmd := exec.Command("mosquitto_sub", "-t", topic)
	outfile, err := os.Create(tempFile)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	time.Sleep(10 * time.Second)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(20 * time.Second)
	//os.Remove(tempFile)

}



