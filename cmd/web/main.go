package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"net/http"
	"os"

	"github.com/talosrobert/golang-boxes/internal/models"
)

type application struct {
	logger        zerolog.Logger
	boxes         *models.BoxModel
	templateCache templateCache
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /box/view/{id}", app.boxView)
	mux.HandleFunc("GET /box/create", app.boxCreate)
	mux.HandleFunc("POST /box/create", app.boxCreatePost)

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	return app.logRequest(mux)
}

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	addr := flag.String("addr", ":4321", "HTTP server network address")
	dsn := flag.String("dsn", "localhost:5432", "database network address")
	flag.Parse()

	dbpool, err := pgxpool.New(context.Background(), *dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to create connection pool")
		os.Exit(1)
	}
	defer dbpool.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to create connection pool")
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		boxes:         &models.BoxModel{DB: dbpool},
		templateCache: templateCache,
	}

	mux := app.routes()

	logger.Info().Msgf("Starting HTTP server on %s", *addr)
	if err = http.ListenAndServe(*addr, mux); err != nil {
		logger.Fatal().Err(err).Send()
		os.Exit(1)
	}
}
