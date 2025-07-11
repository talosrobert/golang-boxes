package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.logger.Error().Err(err).Str("http_method", r.Method).Str("uri", r.URL.RequestURI()).Send()
		http.NotFound(w, r)
		return
	}

	box, err := app.boxes.Get(id)
	if err != nil {
		app.logger.Error().Err(err).Str("http_method", r.Method).Str("uri", r.URL.RequestURI()).Send()
		http.NotFound(w, r)
		return
	}

	var data *templateData
	if flash := app.sessionmanager.PopString(r.Context(), "flash"); flash != "" {
		data = newTemplateData(r, templateDataWithBox(box), templateDataWithFlash(flash))
	} else {
		data = newTemplateData(r, templateDataWithBox(box))
	}

	app.render(w, r, http.StatusOK, "view", data)
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

	app.sessionmanager.Put(r.Context(), "flash", "successfully created a box")
	http.Redirect(w, r, fmt.Sprintf("/box/view/%s", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	var form userCreateForm
	data := newTemplateData(r, templateDataWithForm(form))
	app.render(w, r, http.StatusOK, "signup", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &userCreateForm{
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
		Psw:   r.PostForm.Get("password"),
	}

	form.CheckField("name", "This field cannot be blank", validator.NotBlank(form.Name))
	form.CheckField("email", "This field cannot be blank", validator.NotBlank(form.Email))
	form.CheckField("email", "This field must be a valid email address", validator.ValidEmailAddr(form.Email))
	form.CheckField("password", "This field cannot be blank", validator.NotBlank(form.Psw))
	form.CheckField("password", "This field must have at least 6 characters", validator.MinChars(form.Psw, 6))

	if !form.IsValid() {
		data := newTemplateData(r, templateDataWithForm(form))
		app.render(w, r, http.StatusUnprocessableEntity, "signup", data)
	}

	_, err = app.users.Insert(form.Name, form.Email, form.Psw)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionmanager.Put(r.Context(), "flash", "successfully created a user")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request)      {}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request)  {}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {}
