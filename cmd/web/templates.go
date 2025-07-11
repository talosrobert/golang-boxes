package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/talosrobert/golang-boxes/internal/models"
	"github.com/talosrobert/golang-boxes/internal/validator"
)

type userCreateForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Psw                 string `form:"password"`
	validator.Validator `form:"-"`
}

type boxCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type templateData struct {
	CurrentYear int
	Box         models.Box
	Boxes       []models.Box
	Form        any
	Flash       string
}

func newTemplateData(r *http.Request, opts ...func(*templateData)) *templateData {
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

func templateDataWithForm(form any) func(*templateData) {
	return func(td *templateData) {
		td.Form = form
	}
}

func templateDataWithFlash(flash string) func(*templateData) {
	return func(td *templateData) {
		td.Flash = flash
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
