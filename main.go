package main

import (
	"net/http"
	"html/template"
	"path"
	"log"
)

type BasicElements struct {
	Title string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	page := BasicElements{"Index"}
	fp := path.Join("html/templates/", "index.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(w, page); err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)

}
