package main

import (
	"net/http"

	"github.com/apigban/lenslocked_v1/views"
	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil) //w - where the template outputs to, homeView.Layout - Template/layout that needs to be rendered, nil - data to be passed
		if err != nil {
		panic(err) //TODO - handle error instead of panicking
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Template.ExecuteTemplate(w, contactView.Layout, nil) //w - where the template outputs to, homeView.Layout - Template/layout that needs to be rendered, nil - data to be passed
	if err != nil {
		panic(err) //TODO - handle error instead of panicking
	}
	}
}

func main() {

	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
