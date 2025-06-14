package main

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/talosrobert/golang-boxes/internal/models"
)

type templateData struct {
	Box   models.Box
	Boxes []models.Box
}

type templateCache map[string]*template.Template

func newTemplateCache() (templateCache, error) {
	cache := templateCache{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		tmpls, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		tmpls, err = tmpls.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		tmpls, err = tmpls.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		cache[name] = tmpls
	}

	return cache, nil
}
