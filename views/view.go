package views

import (
	"html/template"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template
	Layout   string
}

// NewView function parses all templates and returns a View type
// Panics when a template cannot be used.
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil { // Parse a view that is not present, will kill the app (panic)
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// layoutFiles returns a slice of strings representing the layout files used by templates
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}
