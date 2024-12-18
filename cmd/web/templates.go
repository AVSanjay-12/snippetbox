package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/AVSanjay-12/snippetbox/internal/models"
)


type templateData struct{
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken string
}

func humanDate(t time.Time) string{
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")

}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error){
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil{
		return nil, err
	}

	for _, page := range pages{
		name := filepath.Base(page)
		// The template.FuncMap must be registered with the ts
		// First Parse the base template file into a template set (ts)
		// Add partials to ts
		// Add page template to ts
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil{
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil{
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil{
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}