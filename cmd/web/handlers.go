package main

import (
	"fmt"
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

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	boxes, err := app.boxes.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := newTemplateData(
		r,
		templateDataWithBoxes(boxes),
	)
	app.render(w, r, http.StatusOK, "home", data)
}

func (app *application) boxView(w http.ResponseWriter, r *http.Request) {
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

	data := newTemplateData(
		r,
		templateDataWithBox(box),
	)
	app.render(w, r, http.StatusOK, "view", data)
}

func (app *application) boxCreate(w http.ResponseWriter, r *http.Request) {
	data := newTemplateData(r)
	app.render(w, r, http.StatusOK, "create", data)
}

func (app *application) boxCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.boxes.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/box/view/%d", id), http.StatusSeeOther)
}
