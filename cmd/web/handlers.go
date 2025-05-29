package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error().Err(err).Str("http_method", method).Str("uri", uri).Send()
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	tmpls := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(tmpls...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) boxView(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.logger.Error().Err(err).Str("http_method", r.Method).Str("uri", r.URL.RequestURI()).Send()
		http.NotFound(w, r)
		return
	}

	if id < 1 {
		app.logger.Error().Str("http_method", r.Method).Str("uri", r.URL.RequestURI()).Msg("invalid box id")
		http.NotFound(w, r)
		return
	}

	box, err := app.boxes.Get(id)
	if err != nil {
		app.logger.Error().Err(err).Str("http_method", r.Method).Str("uri", r.URL.RequestURI()).Send()
		http.NotFound(w, r)
		return
	}

	tmpls := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	ts, err := template.ParseFiles(tmpls...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", box)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) boxCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("box create form"))
}

func (app *application) boxCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	rows, err := app.boxes.Insert("this", "sucks", 3)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Write([]byte(fmt.Sprintf("many rows were touched today: %d", rows)))
}
