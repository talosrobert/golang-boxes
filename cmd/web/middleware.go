package main

import (
	"net/http"
)

func (app *application) getSessionId(r *http.Request) string {
	var sessionId string

	cookie, err := r.Cookie(app.sessionmanager.Cookie.Name)
	if err != nil {
		sessionId = "none"
	} else {
		sessionId = cookie.Value
	}

	return sessionId
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip        = r.RemoteAddr
			proto     = r.Proto
			method    = r.Method
			uri       = r.URL.RequestURI()
			sessionId = app.getSessionId(r)
		)

		app.logger.Info().Str("http_proto", proto).Str("client_ip", ip).Str("http_method", method).Str("http_uri", uri).Str("session_id", sessionId).Send()
		next.ServeHTTP(w, r)
	})
}
