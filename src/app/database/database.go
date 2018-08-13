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
	officeLamp := Lights{"Office Lamp", "office_lamp", "true", "officeLamp"}
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
}

func CreateDB() {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const lightsTable = `CREATE TABLE IF NOT EXISTS 
			     lights(id INTEGER PRIMARY KEY, 
			     name TEXT, state TEXT)`

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
	//TODO check if database exists or not

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const insertLight = "INSERT INTO lights (name, state) VALUES (? , ?)"
	const insertLedstrip = `INSERT INTO ledstrips (displayname, name,
				state, color, topic) VALUES (?, ?, ?, ?, ?)`

	for _, light := range InsertKnownLights() {
		lightStatement, _ := db.Prepare(insertLight)
		lightStatement.Exec(light.Name, light.State)
	}

	for _, ledstrip := range InsertKnownLedstrips() {
		ledstripStatement, _ := db.Prepare(insertLedstrip)
		ledstripStatement.Exec(ledstrip.DisplayName, ledstrip.Name,
				       ledstrip.State, ledstrip.Color,
				       ledstrip.Topic)
	}
}

func GetlightState() {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const getLightState = "SELECT name, state FROM lights"

	//getLighStateStatement, err := db.Prepare(getLightState)
	rows, err := db.Query(getLightState)
	if err != nil {
		log.Fatal(err)
	}

	var name  string
	var state string

	for rows.Next() {
		rows.Scan(&name, &state)
	}

	//for rows.Next() {
	//	rows.Scan(&state)
	//}
	//fmt.Println("State: ", state)



	//getLighStateStatement.Exec(light.Name)

	//return currentState

	//for rows.Next() {
	//	rows.Scan(&name, &state)
	//}

}

func UpdateLights(light Lights, state string) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const updateLight = "UPDATE lights SET state = ? WHERE NAME = ?"

	updateLightStatement, err := db.Prepare(updateLight)
	if err != nil {
		log.Fatal(err)
	}
	updateLightStatement.Exec(state, light.Name)
}

func UpdateLedstrip(name string, color string, state string) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	const updateLedstrip = `UPDATE ledstrips SET state = ?, color = ?
	                        WHERE name = ?`
	updateLedstripStatement, err := db.Prepare(updateLedstrip)
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
