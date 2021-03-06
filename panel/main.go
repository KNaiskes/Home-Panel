package main

import (
	"html/template"
	"log"
	"net/http"
	"github.com/KNaiskes/Home-Panel/mqtt"
	"github.com/KNaiskes/Home-Panel/database"
	"github.com/gorilla/sessions"
)

var htmlTemplates = "src/github.com/KNaiskes/Home-Panel/html/templates/"

type LoginMessages struct {
	CredentialsMatch bool
	UsernameLength int
	PasswordLength int
}

var store = sessions.NewCookieStore([]byte("keep-it-safe-keep-it-secret"))
var userIsAdmin bool

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/ledstrip", ledStripHandler)
	http.HandleFunc("/twoState", twoStateDevicesHandler)
	http.HandleFunc("/admin-panel", adminPanelHandler)
	http.HandleFunc("/addUser", addUserHander)
	http.HandleFunc("/delUser", deleteUserHandler)
	http.HandleFunc("/updatePass", updatePassHandler)
	http.HandleFunc("/addNewDevice", addNewDeviceHandler)
	http.HandleFunc("/removeDevice", removeTwoStateDeviceHandler)
	http.HandleFunc("/temperatureHum", tempHumHandler)
	http.Handle("/src/github.com/KNaiskes/Home-Panel/html/static/",
		http.StripPrefix("/src/github.com/KNaiskes/Home-Panel/html/static/",
		http.FileServer(http.Dir("src/github.com/KNaiskes/Home-Panel/html/static/"))))

	http.ListenAndServe(":8080", nil)

}

func isLoggedIn(sessionName string, w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		//http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
}

func isAdmin(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal(err)
	}
	if session.Values["username"] != "admin" {
		http.Redirect(w, r, "dashboard", http.StatusSeeOther)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fp := htmlTemplates + "index.html"
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

	lenUsername := len(usernameForm)
	lenPassword := len(passwordForm)
	match := database.CheckUser(usernameForm, passwordForm)

	Messages := LoginMessages{match, lenUsername, lenPassword}

	if database.CheckUser(usernameForm, passwordForm) == true {
		session, err := store.Get(r, "cookie-name")
		if err != nil {
			log.Fatal(err)
		}
		session.Values["authenticated"] = true
		session.Values["username"] = usernameForm
		session.Save(r, w)

		if usernameForm == "admin" && passwordForm == "admin" {
			http.Redirect(w, r, "updatePass", http.StatusSeeOther)
		}

		if usernameForm == "admin" {
			userIsAdmin = true
		} else {
			userIsAdmin = false
		}

		http.Redirect(w, r, "dashboard", http.StatusSeeOther)
	}

	fp := htmlTemplates + "login.html"

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, Messages)
	if err != nil {
		log.Fatal(err)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)

	fp := htmlTemplates + "dashboard.html"

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, userIsAdmin)
	if err != nil {
		log.Fatal(err)
	}

}

func ledStripHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)

	fp := htmlTemplates + "ledstrip.html"

	for _, ledstrip := range database.DBledstrips() {
		ledstrip_state := r.FormValue(ledstrip.Name)
		ledstrip_color := r.FormValue(ledstrip.Color)

		if ledstrip_state == "true" {
			mqtt.ChangeState("on", ledstrip.Topic)
		} else if ledstrip_state == "false" {
			mqtt.ChangeState("off", ledstrip.Topic)
		}

		if ledstrip_state == "" {
			ledstrip_state = ledstrip.State
		}
		if ledstrip_color != "" {
			database.UpdateLedstrip(ledstrip.Name, ledstrip_color,
				ledstrip_state)
			mqtt.ChangeColor(ledstrip_color, ledstrip.Topic)
		}

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

func twoStateDevicesHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)

	fp := htmlTemplates + "twoState.html"
	tmpl, err := template.ParseFiles(fp)

	for _, device := range database.DBtwoStateDevices() {
		device_state := r.FormValue(device.Name)

		if device_state == "true" {
			database.UpdateTwoState(device.Name, device_state)
		} else if device_state == "false" {
			database.UpdateTwoState(device.Name, device_state)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, database.DBtwoStateDevices())
	if err != nil {
		log.Fatal(err)
	}
}

func tempHumHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)

	fp := htmlTemplates + "tempHumTable.html"
	tmpl, err := template.ParseFiles(fp)

	metrics := database.GetTempHum()

	var temperature     []string
	var humidity        []string
	var time            []string

	for  t := 0; t < len(metrics); t += 3 {
		temperature = append(temperature, metrics[t])
	}
	for h := 1; h < len(metrics); h += 3 {
		humidity = append(humidity, metrics[h])
	}
	for d := 2; d < len(metrics); d += 3 {
		time = append(time, metrics[d])
	}

	measurements := Measurements{temperature, humidity, time}

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, measurements)
	if err != nil {
		log.Fatal(err)
	}
}
