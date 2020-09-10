package handlers

import (
	"log"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/rumsrami/pgservice/internal/web"
)

// API ...
func API(build string, shutdown chan os.Signal, database *sqlx.DB, appLogger *log.Logger) *web.App {
	// Overwrite the middleware logger to use the app logger
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: appLogger, NoColor: false})

	// create a new mux
	app := web.NewApp(shutdown)

	app.Mux.Use(middleware.RequestID)
	app.Mux.Use(middleware.Logger)

	app.Mux.Get("/get-data/{title}", getTitle(database))
	app.Mux.Post("/post-data/title", postTitle(database))

	return app
}
