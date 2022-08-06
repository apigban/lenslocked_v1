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
	homeTemplate, err = template.ParseFiles( // use the global var
		"views/home.gohtml",
		"views/layouts/footer.gohtml",
	)

	if err != nil {
		panic(err) //TODO - handle error instead of panicking
	}

	contactTemplate, err = template.ParseFiles( // use the global var
		"views/contact.gohtml",
		"views/layouts/footer.gohtml",
	)
	if err != nil {
		panic(err) //TODO - handle error instead of panicking
	}

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
