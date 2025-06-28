package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/talosrobert/golang-boxes/internal/models"
)

type application struct {
	logger         zerolog.Logger
	boxes          *models.BoxModel
	templateCache  templateCache
	sessionmanager *scs.SessionManager
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", app.sessionmanager.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("GET /box/view/{id}", app.sessionmanager.LoadAndSave(http.HandlerFunc(app.boxView)))
	mux.Handle("GET /box/create", app.sessionmanager.LoadAndSave(http.HandlerFunc(app.boxCreate)))
	mux.Handle("POST /box/create", app.sessionmanager.LoadAndSave(http.HandlerFunc(app.boxCreatePost)))

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
		logger:         logger,
		boxes:          &models.BoxModel{DB: dbpool},
		templateCache:  templateCache,
		sessionmanager: scs.New(),
	}

	app.sessionmanager.Store = pgxstore.New(dbpool)
	app.sessionmanager.Lifetime = (time.Hour * 12)
	mux := app.routes()

	logger.Info().Msgf("Starting HTTP server on %s", *addr)
	if err = http.ListenAndServe(*addr, mux); err != nil {
		logger.Fatal().Err(err).Send()
		os.Exit(1)
	}
}
