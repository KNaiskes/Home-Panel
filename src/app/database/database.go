package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"app/devices"
	"log"
)

func CreateDB() {
	db, err := sql.Open("sqlite3", "home.db")
	if err != nil {
		log.Fatal(err)
	}

	const lightsTable = `CREATE TABLE IF NOT EXISTS 
			     lights(id INTEGER PRIMARY KEY, 
			     name TEXT, state TEXT)`

	const ledstripsTable = `CREATE TABLE IF NOT EXISTS
			       ledstrips(id INTEGER PRIMARY KEY,
			       name TEXT, color TEXT, state TEXT)`

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

	db, err := sql.Open("sqlite3", "home.db")
	if err != nil {
		log.Fatal(err)
	}

	const insertLight = "INSERT INTO lights (name, state) VALUES (? , ?)"
	const insertLedstrip = "INSERT INTO ledstrips (name, color, state) VALUES (?, ?,?)"

	for _, light := range devices.GetLights() {
		lightStatement, _ := db.Prepare(insertLight)
		lightStatement.Exec(light.Name, light.State)
	}

	for _, ledstrip := range devices.GetLedstrips() {
		ledstripStatement, _ := db.Prepare(insertLedstrip)
		ledstripStatement.Exec(ledstrip.Name, ledstrip.Color, ledstrip.State)
	}
}
