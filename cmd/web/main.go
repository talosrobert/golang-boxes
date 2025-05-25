package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

type application struct {
	logger zerolog.Logger
}

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /box/view/{id}", app.boxView)
	mux.HandleFunc("GET /box/create", app.boxCreate)
	mux.HandleFunc("POST /box/create", app.boxCreatePost)

	return mux
}

func main() {
	app := &application{
		logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	addr := flag.String("addr", ":4321", "HTTP server network address")
	flag.Parse()

	mux := app.routes()
	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	app.logger.Info().Msgf("Starting HTTP server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	app.logger.Fatal().Msg(err.Error())
	os.Exit(1)
}
