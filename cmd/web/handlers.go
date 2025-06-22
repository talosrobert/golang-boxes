package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/talosrobert/golang-boxes/internal/validator"
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

type boxCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *application) boxCreate(w http.ResponseWriter, r *http.Request) {
	form := boxCreateForm{Expires: 365}
	data := newTemplateData(r, templateDataWithForm(form))
	app.render(w, r, http.StatusOK, "create", data)
}

func (app *application) boxCreatePost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &boxCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	form.CheckField("title", "this field cannot be left blank", validator.NotBlank(form.Title))
	form.CheckField("title", "this field cannot be longer then 100 chars", validator.MaxChars(form.Title, 100))
	form.CheckField("content", "this field cannot be left blank", validator.NotBlank(form.Content))
	form.CheckField("expires", "this field must equal 1, 7 or 365", validator.PermittedValues(form.Expires, 1, 7, 365))

	if !form.IsValid() {
		app.logger.Error().Msg("invalid user input in boxCreateForm")
		data := newTemplateData(
			r,
			templateDataWithForm(form),
		)
		app.render(w, r, http.StatusUnprocessableEntity, "create", data)
		return
	}

	id, err := app.boxes.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/box/view/%d", id), http.StatusSeeOther)
}
