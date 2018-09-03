package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"fmt"
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

const dbDir = "src/app/db/"
const dbName = dbDir + "home.db"
const dbUsers = dbDir + "users.db"

func CreateUsersDB() {
	db, err := sql.Open("sqlite3", dbUsers)
	if err != nil {
		log.Fatal(err)
	}

	const userTable = `CREATE TABLE IF NOT EXISTS
			   users(id INTEGER PRIMARY KEY, username TEXT,
			   password TEXT)`
	statement, err := db.Prepare(userTable)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func AddUser(username string, password string) {
	db, err := sql.Open("sqlite3", dbUsers)
	if err != nil {
		log.Fatal(err)
	}
	const insertUser = `INSERT INTO users(username, password) VALUES (?, ?)`
	statement, err := db.Prepare(insertUser)
	statement.Exec(username, password)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckUser(username string, password string) bool {
	db, err := sql.Open("sqlite3", dbUsers)
	if err != nil {
		log.Fatal(err)
	}
	statement := `SELECT username, password  FROM users WHERE username=? AND password=?`
	err = db.QueryRow(statement, username, password).Scan(&username, &password)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return false
	}
	return true
}

func UserExists(username string) bool {
	db, err := sql.Open("sqlite3", dbUsers)
	if err != nil {
		log.Fatal(err)
	}
	statement := `SELECT username FROM users WHERE username=?`
	err = db.QueryRow(statement, username).Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return false
	}
	return true
}

func InsertKnownLedstrips() []LedStrip {
	//Already known led strips that will be added only when database is
	//lost or about to be created
	//Add any new led strips below
	bedroomLedstrip := LedStrip{"Bedroom", "bedroom_ledstrip", "false",
				   "white", "ledStrip"}

	MyledStrips := []LedStrip{bedroomLedstrip}

	return MyledStrips
}

func InsertKnownLights() []Lights {
	//Already known lights that will be added only when database is lost
	//or about to be created
	//Add any new light below
	officeLamp := Lights{"Office Lamp", "office_lamp", "false", "officeLamp"}
	DeskLamp := Lights{"Desk Lamp", "desk_lamp", "false", "deskLamp"}
	MyLights := []Lights{officeLamp, DeskLamp}

	return MyLights
}

func DBexists() {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0700)
		CreateDB()
		InsertAll()
	}
	if _, err := os.Stat(dbUsers); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0700)
		CreateUsersDB()
		AddUser("admin", "admin")
	}

}

func CreateDB() {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const lightsTable = `CREATE TABLE IF NOT EXISTS 
			     lights(id INTEGER PRIMARY KEY, displayname TEXT,
			     name TEXT, state TEXT, topic TEXT)`

	const ledstripsTable = `CREATE TABLE IF NOT EXISTS
			       ledstrips(id INTEGER PRIMARY KEY,
			       displayname TEXT, name TEXT, state TEXT,
			       color TEXT, topic TEXT)`

	statement, err := db.Prepare(lightsTable)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = db.Prepare(ledstripsTable)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func InsertAll() {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const insertLight = `INSERT INTO lights (displayname, name, state, topic) VALUES (?, ?, ?, ?)`
	const insertLedstrip = `INSERT INTO ledstrips (displayname, name,
				state, color, topic) VALUES (?, ?, ?, ?, ?)`

	for _, light := range InsertKnownLights() {
		lightStatement, _ := db.Prepare(insertLight)
		lightStatement.Exec(light.DisplayName, light.Name,
				    light.State, light.Topic)
	}

	for _, ledstrip := range InsertKnownLedstrips() {
		ledstripStatement, _ := db.Prepare(insertLedstrip)
		ledstripStatement.Exec(ledstrip.DisplayName, ledstrip.Name,
				       ledstrip.State, ledstrip.Color,
				       ledstrip.Topic)
	}
}

func DBlights() []Lights {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	var displayname string
	var name string
	var state string
	var topic string

	Mylights := []Lights{}
	const getLightState = `SELECT displayname, name, state, 
			       topic FROM lights`

	//getLighStateStatement, err := db.Prepare(getLightState)
	rows, err := db.Query(getLightState)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&displayname, &name, &state, &topic)
		temp := Lights{displayname, name, state, topic}
		Mylights = append(Mylights, temp)
	}
	return Mylights
}

func UpdateLights(name string, state string) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const updateLight = "UPDATE lights SET state = ? WHERE name = ?"

	updateLightStatement, err := db.Prepare(updateLight)
	if err != nil {
		log.Fatal(err)
	}
	updateLightStatement.Exec(state, name)
}

func UpdateLedstrip(name string, color string, state string) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	const updateLedstrip = `UPDATE ledstrips SET state = ?, color = ?
	                        WHERE name = ?`
	updateLedstripStatement, err := db.Prepare(updateLedstrip)
	if err != nil {
		log.Fatal(err)
	}
	updateLedstripStatement.Exec(state, color, name)
}

func DBledstrips() []LedStrip {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	var displayName string
	var name string
	var state string
	var color string
	var topic string

	MyLedstrips := []LedStrip{}

	const getLedstrips = `SELECT displayname, name, state, color, topic
			     FROM ledstrips`
	rows, err := db.Query(getLedstrips)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&displayName, &name, &state, &color, &topic)
		temp := LedStrip{displayName, name, state, color, topic}
		MyLedstrips = append(MyLedstrips, temp)
	}
	fmt.Println("MyLedstrips:", MyLedstrips)
	return MyLedstrips
}
