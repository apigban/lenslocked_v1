package main

import (
	"net/http"

	"github.com/apigban/lenslocked_v1/controllers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Static Controller
	staticC := controllers.NewStatic()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	//Users controller
	usersC := controllers.NewUsers()
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
