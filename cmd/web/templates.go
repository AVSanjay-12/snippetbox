package main

import (
	"path/filepath"
	"text/template"

	"github.com/AVSanjay-12/snippetbox/internal/models"
)


type templateData struct{
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error){
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil{
		return nil, err
	}

	for _, page := range pages{
		name := filepath.Base(page)
		// First Parse the base template file into a template set (ts)
		// Add partials to ts
		// Add page template to ts
		ts, err := template.ParseFiles("./ui/html/base.html")
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