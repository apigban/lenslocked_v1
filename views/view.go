package views

import (
	"html/template"
)

type View struct {
	Template *template.Template
}

// NewView function parses all templates and returns a View type
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err) // NewView only needs to be used during setup - failure to parse a template kills the app
	}

	return &View{
		Template: t,
	}
}
