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

type TwoState struct {
	DisplayName string
	Name	    string
	State       string
	Topic       string
}

const dbDir = "src/github.com/KNaiskes/Home-Panel/db/"
const dbName = dbDir + "home.db"
const dbUsers = dbDir + "users.db"
const dbMeasurements = dbDir + "measurements.db"

const DriverDB = "sqlite3"

func SimpleQuery(driver string, dbName string, query string) {
	db, err := sql.Open(driver, dbName)
	if err != nil {
		log.Fatal(err)
	}
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func CreateUsersDB() {
	query := `CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY,
			username TEXT, password TEXT)`
	SimpleQuery(DriverDB, dbUsers, query)


}

func CreateMeasurementsDB() {
	query := `CREATE TABLE IF NOT EXISTS 
			measurements(id INTEGER PRIMARY KEY, 
			temperature REAL, humidity REAL)`
	SimpleQuery(DriverDB, dbMeasurements, query)
}

func CreateDB() {
	//TODO: rename function to something more relative
	//TODO: rename lights table to twoState
	//TODO: rename database name to something more relative

	twoStateQuery := `CREATE TABLE IF NOT EXISTS 
			lights(id INTEGER PRIMARY KEY, displayname TEXT, 
			name TEXT, state TEXT, topic TEXT)`
	SimpleQuery(DriverDB, dbName, twoStateQuery)

	ledStripQuery := `CREATE TABLE IF NOT EXISTS 
				  ledstrips(id INTEGER PRIMARY KEY,
				  displayname TEXT, name TEXT, state TEXT,
				  color TEXT, topic TEXT)`
	SimpleQuery(DriverDB, dbName, ledStripQuery)
}

func AddTempHum(temperature float64, humidity float64) {
	db, err := sql.Open(DriverDB, dbMeasurements)
	if err != nil {
		log.Fatal(err)
	}
	const addTemperatureTable = `INSERT INTO measurements(temperature, humidity) VALUES(?, ?)`
	statement, err := db.Prepare(addTemperatureTable)
	statement.Exec(temperature, humidity)
	if err != nil {
		log.Fatal(err)
	}
}

func GetTemperature() []float64 {
	db, err := sql.Open(DriverDB, dbMeasurements)
	if err != nil {
		log.Fatal(err)
	}
	metrics := []float64{}
	var Temperature float64
	const getTempTable = `SELECT temperature FROM measurements`
	rows, err := db.Query(getTempTable)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&Temperature)
		metrics = append(metrics, Temperature )
	}
	return metrics
}

func GetHumidity() []float64 {
	db, err := sql.Open(DriverDB, dbMeasurements)
	if err != nil {
		log.Fatal(err)
	}
	metrics := []float64{}
	var Humidity float64
	const getHumTable = `SELECT temperature FROM measurements`
	rows, err := db.Query(getHumTable)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&Humidity)
		metrics = append(metrics, Humidity )
	}
	return metrics
}

func AddUser(username string, password string) bool {
	db, err := sql.Open(DriverDB, dbUsers)
	if err != nil {
		log.Fatal(err)
	}
	/*
	const userExists = `SELECT username FROM users WHERE username=?`

	err = db.QueryRow(userExists, username).Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return false
	}
	*/
	if !UserExists(username) {
		if len(username) >= 5 && len(password) >= 5 {
			const insertUser = `INSERT INTO users(username, password) VALUES (?, ?)`
			statement, err := db.Prepare(insertUser)
			statement.Exec(username, password)
			if err != nil {
				log.Fatal(err)
			}
		}
		return true
	}
	return false
}

func CheckUser(username string, password string) bool {
	db, err := sql.Open(DriverDB, dbUsers)
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
	db, err := sql.Open(DriverDB, dbUsers)
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


func DelUser(username string) bool {
	db, err := sql.Open(DriverDB, dbUsers)
	if err != nil {
		log.Fatal(err)
	}
	const userExists = `SELECT username FROM users WHERE username=?`
	const delUser = `DELETE FROM users WHERE username=?`

	err = db.QueryRow(userExists, username).Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return false
	}
	statement, err := db.Prepare(delUser)
	statement.Exec(&username)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func ShowUsers() []string {
	db, err := sql.Open(DriverDB, dbUsers)
	if err != nil {
		log.Fatal(err)
	}
	var username string
	usernames := []string{}

	const showUsers = `SELECT username FROM users`

	rows, err := db.Query(showUsers)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&username)
		usernames = append(usernames, username)
	}
	return usernames
}

func UpdatePassword(username string, password string) {
	db, err := sql.Open(DriverDB, dbUsers)
	if err != nil {
		log.Fatal(err)
	}

	const updatePass = `UPDATE users set password=? WHERE username=?`
	statement, err := db.Prepare(updatePass)
	statement.Exec(password, username)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertKnownLedstrips() []LedStrip {
	bedroomLedstrip := LedStrip{"Bedroom", "bedroom_ledstrip", "false",
				   "white", "ledStrip"}

	MyledStrips := []LedStrip{bedroomLedstrip}

	return MyledStrips
}

func InsertKnownDevices() []TwoState {
	officeLamp := TwoState{"Office Lamp", "office_lamp", "false", "officeLamp"}
	DeskLamp := TwoState{"Desk Lamp", "desk_lamp", "false", "deskLamp"}
	MyDevices := []TwoState{officeLamp, DeskLamp}

	return MyDevices
}

func DeviceExists(displayname string, name string, topic string) bool {
	db, err := sql.Open(DriverDB, dbName)
	if err != nil {
		log.Fatal(err)
	}
	statement := `SELECT displayname, name, topic FROM lights 
			WHERE displayname=? OR name=? OR topic=?`
	err = db.QueryRow(statement, displayname, name, topic).Scan(&displayname, &name, &topic)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return false
	}
	return true
}

func AddTwoStateDevice(displayname string, name string, topic string) bool {
	db, err := sql.Open(DriverDB, dbName)
	if err != nil {
		log.Fatal(err)
	}
	if !DeviceExists(displayname, name, topic) && len(name) >= 3 && len(displayname) >= 3 && len(topic) >= 3 {
		state := "false"
		const checkDevice = `INSERT INTO lights (displayname, name, state, topic) VALUES (?, ?, ?, ?)`
		statement, err := db.Prepare(checkDevice)
		statement.Exec(displayname, name, state, topic)

		if err != nil {
			log.Fatal(err)
		}
		return true
	}
	return false
}

func RemoveTwoStateDevice(name string) bool {
	db, err := sql.Open(DriverDB, dbName)
	if err != nil {
		log.Fatal(err)
	}
	const deviceExists = `SELECT name FROM lights WHERE name=?`
	const deleteDevice = `DELETE FROM lights WHERE name=?`

	err = db.QueryRow(deviceExists, name).Scan(&name)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return false
	}
	statement, err := db.Prepare(deleteDevice)
	statement.Exec(&name)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func AvailableDevices() []string {
	db, err := sql.Open(DriverDB, dbName)
	if err != nil {
		log.Fatal(err)
	}
	var deviceName string
	devices := []string{}

	const getDevices = `SELECT name FROM lights`

	rows, err := db.Query(getDevices)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&deviceName)
		devices = append(devices, deviceName)
	}
	return devices
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
	if _, err := os.Stat(dbMeasurements); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0700)
		CreateMeasurementsDB()
	}

}


func InsertAll() {
	db, err := sql.Open(DriverDB, dbName)
	if err != nil {
		log.Fatal(err)
	}

	const insertDevice = `INSERT INTO lights (displayname, name, state, topic) VALUES (?, ?, ?, ?)`
	const insertLedstrip = `INSERT INTO ledstrips (displayname, name,
				state, color, topic) VALUES (?, ?, ?, ?, ?)`

	for _, device := range InsertKnownDevices() {
		deviceStatement, _ := db.Prepare(insertDevice)
		deviceStatement.Exec(device.DisplayName, device.Name,
				    device.State, device.Topic)
	}

	for _, ledstrip := range InsertKnownLedstrips() {
		ledstripStatement, _ := db.Prepare(insertLedstrip)
		ledstripStatement.Exec(ledstrip.DisplayName, ledstrip.Name,
				       ledstrip.State, ledstrip.Color,
				       ledstrip.Topic)
	}
}

func DBtwoStateDevices() []TwoState{
	db, err := sql.Open(DriverDB, dbName)
	if err != nil {
		log.Fatal(err)
	}

	var displayname string
	var name string
	var state string
	var topic string

	TwoStateDevices := []TwoState{}
	const getDeviceState = `SELECT displayname, name, state, 
			       topic FROM lights`

	rows, err := db.Query(getDeviceState)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&displayname, &name, &state, &topic)
		temp := TwoState{displayname, name, state, topic}
		TwoStateDevices = append(TwoStateDevices, temp)
	}
	return TwoStateDevices
}

func UpdateTwoState(name string, state string) {
	db, err := sql.Open(DriverDB, dbName)
	if err != nil {
		log.Fatal(err)
	}

	const updateState = "UPDATE lights SET state = ? WHERE name = ?"

	updateStateStatement, err := db.Prepare(updateState)
	if err != nil {
		log.Fatal(err)
	}
	updateStateStatement.Exec(state, name)
}

func UpdateLedstrip(name string, color string, state string) {
	db, err := sql.Open(DriverDB, dbName)
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
	db, err := sql.Open(DriverDB, dbName)
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
