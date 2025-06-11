package main

import (
	"html/template"
	"path/filepath"

	"github.com/talosrobert/golang-boxes/internal/models"
)

type templateData struct {
	Box   models.Box
	Boxes []models.Box
}

type templateCache = map[string]*template.Template

func newTemplateCache() (templateCache, error) {
	cache := templateCache{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
