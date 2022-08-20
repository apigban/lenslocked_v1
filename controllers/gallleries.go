package controllers

import (
	"github.com/apigban/lenslocked_v1/models"
	"github.com/apigban/lenslocked_v1/views"
)

//
// TODO - SET method GET /
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  &gs,
	}
}

type Galleries struct {
	New *views.View
	gs  *models.GalleryService
}
