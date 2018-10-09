package main

import (
	"net/http"
	"log"
	"html/template"
	"github.com/KNaiskes/Home-Panel/database"
)

type AddUserMessages struct {
	UsernameLength int
	PasswordLength int
	AddedUser bool
}

type UpdatePasswordMessages struct {
	UsernameExists bool
	UserNameLength int
	NewPasswordLength int
	UsernamesList []string
}

type DelUserMessages struct {
	UsernameLength int
	UsernameExists bool
	UsernamesList []string
}

type Measurements struct {
	Temperature []string
	Humidity    []string
	Time        []string
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

	fp := htmlTemplates + "adminPanel.html"
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
	userAdded := database.AddUser(registerUsername, registerPassword)

	Messages := AddUserMessages{lenUsername, lenPassword, userAdded}

	fp := htmlTemplates + "addUser.html"
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, Messages)
	if err != nil {
		log.Fatal(err)
	}

	//database.AddUser(registerUsername, registerPassword)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	delUsernameForm := r.FormValue("username")

	lenUsername := len(delUsernameForm)
	UserExists := database.UserExists(delUsernameForm)
	usersList := database.ShowUsers()

	Messages := DelUserMessages{lenUsername, UserExists, usersList}

	fp := htmlTemplates + "delUser.html"
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

func addNewDeviceHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	name := r.FormValue("deviceName")
	dispayName := r.FormValue("displayName")
	mqttTopic := r.FormValue("deviceMqtt")

	fp := htmlTemplates + "addNewDevice.html"
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

	fp := htmlTemplates + "updatePass.html"
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

func removeTwoStateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn("cookie-name", w, r)
	isAdmin(w, r)

	name := r.FormValue("deviceName")

	fp := htmlTemplates + "removeTwoStateDevice.html"
	tmpl, err := template.ParseFiles(fp)

	database.RemoveTwoStateDevice(name)

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, database.AvailableDevices())
	if err != nil {
		log.Fatal(err)
	}
}


