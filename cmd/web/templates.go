package main

import (
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/talosrobert/golang-boxes/internal/models"
)

type templateData struct {
	CurrentYear int
	Box         models.Box
	Boxes       []models.Box
}

func newTemplateData(opts ...func(*templateData)) *templateData {
	td := &templateData{
		CurrentYear: time.Now().Year(),
	}

	for _, o := range opts {
		o(td)
	}

	return td
}

func templateDataWithBox(box models.Box) func(*templateData) {
	return func(td *templateData) {
		td.Box = box
	}
}

func templateDataWithBoxes(boxes []models.Box) func(*templateData) {
	return func(td *templateData) {
		td.Boxes = boxes
	}
}

type templateCache map[string]*template.Template

func newTemplateCache() (templateCache, error) {
	cache := templateCache{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		cache[name] = ts
	}

	return cache, nil
}
