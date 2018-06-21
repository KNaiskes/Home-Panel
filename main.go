package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It works!")
}

func main() {

	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)

}
