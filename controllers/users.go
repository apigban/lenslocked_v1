package controllers

import (
	"fmt"
	"net/http"

	"github.com/apigban/lenslocked_v1/models"
	"github.com/apigban/lenslocked_v1/views"
)

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
//
// GET /signup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is used to process the signup form when a user submits it.
// this is used to create a new user account
//
// POST /signup
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

// Login is used to verify the user provided user and password
// Allows user to access if information are correct.
//
// POST /login
func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	// user := models.User{
	// 	Email:    form.Email,
	// 	Password: form.Password,
	// }

	// // if err := u.us.Create(&user); err != nil {
	// // 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// // 	return
	// // }
	fmt.Fprintln(w, form)
}
