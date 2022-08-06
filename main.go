package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Global var to try out using html/templates
var (
	homeTemplate    *template.Template
	contactTemplate *template.Template
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		panic(err) //TODO - handle error instead of panicking
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func main() {
	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml") // use the global var
	if err != nil {
		panic(err) //TODO - handle error instead of panicking
	}

	contactTemplate, err = template.ParseFiles("views/contact.gohtml")
	if err != nil {
		panic(err) //TODO - handle error instead of panicking
	}

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
