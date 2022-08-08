package controllers

import "github.com/apigban/lenslocked_v1/views"

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/home.gohtml"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
