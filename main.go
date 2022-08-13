package main

import (
	"fmt"
	"net/http"

	"github.com/apigban/lenslocked_v1/controllers"
	"github.com/apigban/lenslocked_v1/models"
	"github.com/gorilla/mux"
)

// TODO - Fix before prod
const (
	host     = "localhost"
	port     = 5432
	user     = "PGUSER"
	password = "PASSWORD"
	dbname   = "lenslocked_test"
)

func main() {
	// TODO - Fix before prod
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	must(err)
	defer us.Close()
	us.AutoMigrate()
	// us.DestructiveReset()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
