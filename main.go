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
	err = tmpl.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.HandleFunc("/", indexHandler)
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/html/static/", http.StripPrefix("/html/static/", http.FileServer(http.Dir("html/static/"))))

	http.ListenAndServe(":8080", nil)

}
