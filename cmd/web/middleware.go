package main

import (
	"net/http"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info().Str("http_proto", proto).Str("client_ip", ip).Str("http_method", method).Str("http_uri", uri).Send()

		next.ServeHTTP(w, r)
	})
}
