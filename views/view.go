package views

import (
	"html/template"
)

type View struct {
	Template *template.Template
}

// NewView function parses all templates and returns a View type
// Panics when a template cannot be used.
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil { // Parse a view that is not present, will kill the app (panic)
		panic(err)
	}

	return &View{
		Template: t,
	}
}
