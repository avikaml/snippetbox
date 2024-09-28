package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
	"github.com/avikaml/snippetbox/internal/models"
	"github.com/avikaml/snippetbox/ui" 
)
type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken string
}

func humanDate(t time.Time) string {
	if t.IsZero(){
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate" : humanDate,
}


func newTemplateCache() (map[string]*template.Template, error){
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil{
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		// if err != nil{
		// 	return nil, err
		// }

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.html",
			page,
		}


		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		// files := []string{
		// 	"./ui/html/base.tmpl.html",
		// 	"./ui/html/partials/nav.tmpl.html",
		// 	page,
		// }

		// ts, err := template.ParseFiles(files...)
		// ts, err = ts.ParseFiles(page)
		// if err != nil{
		// 	return nil, err
		// }
		cache[name] = ts
	}
	return cache, nil
}

