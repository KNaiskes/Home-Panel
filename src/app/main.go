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

type AddUserMessages struct {
	UsernameLength int
	PasswordLength int
	UsernameExists bool
}

type UpdatePasswordMessages struct {
	UsernameExists bool
	UserNameLength int
	NewPasswordLength int
	UsernamesList []string
}

type LoginMessages struct {
	CredentialsMatch bool
	UsernameLength int
	PasswordLength int
}

type DelUserMessages struct {
	UsernameLength int
	UsernameExists bool
	UsernamesList []string
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
	http.Handle("/src/app/html/static/", http.StripPrefix("/src/app/html/static/",
		http.FileServer(http.Dir("src/app/html/static/"))))

	http.ListenAndServe(":8080", nil)

}

func isLoggedIn(sessionName string, w http.ResponseWriter, r *http.Request) {
	// if user is logged in - create their session
	// else redirect them to log in page

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
	// allow access to some pages only to admin
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal(err)
	}
	if session.Values["username"] != "admin" {
		http.Redirect(w, r, "dashboard", http.StatusSeeOther)
	}
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

	fp := "src/app/html/templates/login.html"

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

	fp := "src/app/html/templates/dashboard.html"

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

func twoStateDevicesHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)

	fp := "src/app/html/templates/twoState.html"
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

func adminPanelHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	//session, err := store.Get(r, "cookie-name")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if session.Values["username"] != "admin" {
	//	http.Redirect(w, r, "dashboard", http.StatusSeeOther)
	//}

	fp := "src/app/html/templates/adminPanel.html"
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func addUserHander(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	registerUsername := r.FormValue("username")
	registerPassword := r.FormValue("password")

	lenUsername := len(registerUsername)
	lenPassword := len(registerPassword)
	userExists  := database.UserExists(registerUsername)

	Messages := AddUserMessages{lenUsername, lenPassword, userExists}

	fp := "src/app/html/templates/addUser.html"
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, Messages)
	if err != nil {
		log.Fatal(err)
	}


	if !database.UserExists(registerUsername) {
		database.AddUser(registerUsername, registerPassword)
	}
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	delUsernameForm := r.FormValue("username")

	lenUsername := len(delUsernameForm)
	UserExists := database.UserExists(delUsernameForm)
	usersList := database.ShowUsers()

	Messages := DelUserMessages{lenUsername, UserExists, usersList}

	fp := "src/app/html/templates/delUser.html"
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, Messages)
	if err != nil {
		log.Fatal(err)
	}

	if delUsernameForm != "admin" && UserExists {
		database.DelUser(delUsernameForm)
	}
}

func updatePassHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	usernameForm := r.FormValue("username")
	passwordForm := r.FormValue("password")
	userExists := database.UserExists(usernameForm)
	lenUsername := len(usernameForm)
	lenPassword := len(passwordForm)
	usersList := database.ShowUsers()

	Messages := UpdatePasswordMessages{userExists, lenUsername, lenPassword, usersList}

	fp := "src/app/html/templates/updatePass.html"
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, Messages)
	if err != nil {
		log.Fatal(err)
	}
	if userExists == true && lenUsername >= 5 && lenPassword >= 5 {
		database.UpdatePassword(usernameForm, passwordForm)
	}
}

func addNewDeviceHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	name := r.FormValue("deviceName")
	dispayName := r.FormValue("displayName")
	mqttTopic := r.FormValue("deviceMqtt")

	exists := database.DeviceExists(name, dispayName, mqttTopic)
	fp := "src/app/html/templates/addNewDevice.html"
	tmpl, err := template.ParseFiles(fp)

	database.AddTwoStateDevice(name, dispayName, mqttTopic)

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
