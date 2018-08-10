package main

import (
	"html/template"
	"log"
	"fmt"
	"net/http"
	"app/mqtt"
)

type LedStrip struct {
	DisplayName string
	Name	    string
	State	    string
	Color       string
	Topic	    string
}

type Lights struct {
	DisplayName string
	Name	    string
	State       string
	Topic       string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/ledstrip", ledStripHandler)
	http.HandleFunc("/lights", lightsHandler)
	http.Handle("/src/app/html/static/", http.StripPrefix("/src/app/html/static/",
		http.FileServer(http.Dir("src/app/html/static/"))))

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

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	fp := "src/app/html/templates/dashboard.html"

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func ledStripHandler(w http.ResponseWriter, r *http.Request) {
	bedroomLedstrip := LedStrip{"Bedroom", "bedroom_ledstrip", "false",
				   "white", "ledStrip"}
	MyledStrips := []LedStrip{bedroomLedstrip}

	fp := "src/app/html/templates/ledstrip.html"

	for _, ledstrip := range MyledStrips {
		ledstrip_state := r.FormValue(ledstrip.Name)
		ledstrip_color := r.FormValue(ledstrip.Color)

		if ledstrip_state == "true" {
			mqtt.ChangeState("on", ledstrip.Topic)
		} else {
			mqtt.ChangeState("off", ledstrip.Topic)
		}
		fmt.Println("State:", ledstrip.State)

		if ledstrip_color != "" {
			mqtt.ChangeColor(ledstrip_color, ledstrip.Topic)
		}
		fmt.Println("Color :", ledstrip_color)

	}

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, MyledStrips)
	if err != nil {
		log.Fatal(err)
	}
}

func lightsHandler(w http.ResponseWriter, r *http.Request) {

	officeLamp := Lights{"Office Lamp", "office_lamp", "true", "officeLamp"}
	DeskLamp := Lights{"Desk Lamp", "desk_lamp", "false", "deskLamp"}
	page := []Lights{officeLamp, DeskLamp}

	fp := "src/app/html/templates/lights.html"
	tmpl, err := template.ParseFiles(fp)

	for _, light := range page {
		light_state := r.FormValue(light.Name)

		if light_state == "true" {
			mqtt.ChangeState("on", light.Topic)
		} else {
			mqtt.ChangeState("off", light.Topic)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

