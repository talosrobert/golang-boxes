package main

import (
	"fmt"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *templateData) {
	tmpls, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("Template %s was not found in templateCache", page)
		app.logger.Error().Err(err).Send()
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	err := tmpls.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.Error().Err(err).Send()
		app.serverError(w, r, err)
		return
	}
}
