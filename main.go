package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

type BasicElements struct {
	Title string
}

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.Handle("/html/static/", http.StripPrefix("/html/static/",
		http.FileServer(http.Dir("html/static/"))))

	http.ListenAndServe(":8080", nil)

}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	page := BasicElements{"Index"}
	fp := path.Join("html/templates/", "index.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	page := BasicElements{"Login-Page"}
	fp := path.Join("html/templates/", "login.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}
