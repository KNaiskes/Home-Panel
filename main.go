package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

type TemperatureHum struct {
	Title string
	DateTime []string
	Temperature []string
}

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/tempHum", temphumiHandler)
	http.Handle("/html/static/", http.StripPrefix("/html/static/",
		http.FileServer(http.Dir("html/static/"))))

	http.ListenAndServe(":8080", nil)

}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("html/templates/", "index.html")

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
	fp := path.Join("html/templates/", "login.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func temphumiHandler(w http.ResponseWriter, r *http.Request) {
	DateTime := []string {"20","29", "40"} //just for testing
	Temperature := []string {"18-06-26", "18-06-26", "18-06-26"} //just for testing
	page := TemperatureHum{"Temperature/Humidity", DateTime, Temperature}
	fp := path.Join("html/templates/", "temphum.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}
