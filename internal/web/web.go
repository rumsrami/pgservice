package web

import (
	"os"
	"syscall"

	"github.com/go-chi/chi"
)

// App is the entry point
type App struct {
	*chi.Mux
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	app := App{
		Mux:      chi.NewRouter(),
		shutdown: shutdown,
	}
	return &app
}

// SignalShutdown is used to gracefully shutdown the app
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
