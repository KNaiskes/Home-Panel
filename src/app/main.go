package main

import (
	"html/template"
	"log"
	"fmt"
	"net/http"
	"app/mqtt"
	"app/database"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("keep-it-safe-keep-it-secret"))

func main() {
	database.DBexists()
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
	usernameForm := r.FormValue("username")
	passwordForm := r.FormValue("password")

	if database.CheckUser(usernameForm, passwordForm) == true {
		session, err := store.Get(r, "cookie-name")
		if err != nil {
			log.Fatal(err)
		}
		session.Values["authenticated"] = true
		session.Save(r, w)
	}

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
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal(err)
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

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

	fp := "src/app/html/templates/ledstrip.html"

	for _, ledstrip := range database.DBledstrips() {
		ledstrip_state := r.FormValue(ledstrip.Name)
		ledstrip_color := r.FormValue(ledstrip.Color)

		if ledstrip_state == "true" {
			//database.UpdateLedstrip(ledstrip.Name, ledstrip_color, "true")
			mqtt.ChangeState("on", ledstrip.Topic)
		} else if ledstrip_state == "false" {
			//database.UpdateLedstrip(ledstrip.Name, ledstrip_color, "false")
			mqtt.ChangeState("off", ledstrip.Topic)
		}
		fmt.Println("State:", ledstrip.State)

		if ledstrip_state == "" {
			ledstrip_state = ledstrip.State
		}
		if ledstrip_color != "" {
			database.UpdateLedstrip(ledstrip.Name, ledstrip_color, ledstrip_state)
			mqtt.ChangeColor(ledstrip_color, ledstrip.Topic)
		}
		fmt.Println("Color :", ledstrip_color)

	}

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, database.DBledstrips())
	if err != nil {
		log.Fatal(err)
	}
}

func lightsHandler(w http.ResponseWriter, r *http.Request) {
	fp := "src/app/html/templates/lights.html"
	tmpl, err := template.ParseFiles(fp)

	for _, light := range database.DBlights() {
		light_state := r.FormValue(light.Name)

		if light_state == "true" {
			database.UpdateLights(light.Name, light_state)
			//mqtt.ChangeState("on", light.Topic)
		} else if light_state == "false" {
			database.UpdateLights(light.Name, light_state)
			//mqtt.ChangeState("off", light.Topic)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, database.DBlights())
	if err != nil {
		log.Fatal(err)
	}
}

