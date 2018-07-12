package main

import (
	"html/template"
	"log"
	"fmt"
	"net/http"
	"os/exec"
	"app/mqtt"
)

func SendMQTT(command string) {
	cmd := exec.Command(command)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed with error: ", err)
	}
	fmt.Println(string(out)) // just for testing
}

type TemperatureHum struct {
	Title string
	DateTime []string
	Temperature []string
}

func main() {
	mqtt.SendMQTT("ls")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/tempHum", temphumiHandler)
	http.Handle("html/static/", http.StripPrefix("html/static/",
		http.FileServer(http.Dir("html/static/"))))

	http.ListenAndServe(":8080", nil)

}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	fp := "src/app/html/templates/index.html"
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fp := "src/app/html/templates/login.html"

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func temphumiHandler(w http.ResponseWriter, r *http.Request) {
	DateTime := []string {"20","29", "40"} //just for testing
	Temperature := []string {"18-06-26", "18-06-26", "18-06-26"} //just for testing
	page := TemperatureHum{"Temperature/Humidity", DateTime, Temperature}
	fp := "src/app/html/templates/temphum.html"

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	fp := "src/app/html/templates/dashboard.html"

	submit := r.FormValue("light1")
	fmt.Println(submit)

	light2 := r.FormValue("light2")
	fmt.Println(light2)


	if submit == "true" {
		fmt.Println("light1 is: true")
	} else {
		fmt.Println("light1 is: false")
	}

	if light2 == "true" {
		fmt.Println("light2 is: true")
	} else {
		fmt.Println("light2 is: false")
	}

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
