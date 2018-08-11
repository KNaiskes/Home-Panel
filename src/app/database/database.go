package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"app/devices"
	"log"
	"os"
)

const dbDir = "src/app/db/"
const dbName = dbDir + "home.db"

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

	db, err := sql.Open("sqlite3", dbName)
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

func GetlightState() (string, string){
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const getLightState = "SELECT name, state FROM lights"
	LightsStates := []devices.Lights

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

func UpdateLights(light devices.Lights, state string) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	const updateLight = "UPDATE lights SET state = ? WHERE NAME = ?"

	updateLightStatement, _ := db.Prepare(updateLight)
	updateLightStatement.Exec(state, light.Name)
}
