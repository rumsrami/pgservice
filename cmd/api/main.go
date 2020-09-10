package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"

	"github.com/rumsrami/pgservice/cmd/api/internal/handlers"
	"github.com/rumsrami/pgservice/internal/db"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"
var serviceName = "API"
var errDatabaseConnection = "Error connecting to db"

func main() {
	if err := run(); err != nil {
		log.Printf("error: %s", err)
	}
	os.Exit(1)
}

func run() error {
	// =========================================================================
	// Logging

	// Set Logging
	infolog := log.New(os.Stdout, fmt.Sprintf("%v Version: %v - ", serviceName, build), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	errlog := log.New(os.Stdout, fmt.Sprintf("%v Version: %v ERROR: - ", serviceName, build), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration

	var cfg struct {
		conf.Version
		Args conf.Args
		Web  struct {
			APIHost         string        `conf:"default:0.0.0.0:5000"`
			ReadTimeout     time.Duration `conf:"default:7s"`
			WriteTimeout    time.Duration `conf:"default:7s"`
			ShutdownTimeout time.Duration `conf:"default:7s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:postgres-db"` // docker-compose database service name
			Name       string `conf:"default:apimain"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	const prefix = "API"
	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// Start Database

	infolog.Println("main : Started : Initializing database support")

	// Start a connection to the database, would fail return err but wont panic
	dbConfig := db.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	database, err := db.OpenTCP(dbConfig)
	if err != nil {
		return errors.Wrap(err, errDatabaseConnection)
	}

	// Check the open connection and return error if not available
	// Shut down if connection is not available
	if err = db.Ping(context.Background(), database); err != nil {
		errlog.Println("main : ", errDatabaseConnection)
		return errors.Wrap(err, errDatabaseConnection)
	}
	defer func() {
		infolog.Printf("main : Database Stopping\n")
		database.Close()
	}()

	// =========================================================================
	// Start Server

	infolog.Println("main : Started : Initializing HTTP support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	// create a new http server with the mux as handler
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      handlers.API(build, shutdown, database, infolog),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// channel to listen for server errors to trigger shutdown
	// Start the service listening for requests.
	go func() {
		infolog.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		infolog.Printf("main : %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			infolog.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}
	return nil
}
