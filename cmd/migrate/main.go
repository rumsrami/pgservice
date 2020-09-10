package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"

	commands "github.com/rumsrami/pgservice/cmd/migrate/command"
	"github.com/rumsrami/pgservice/internal/db"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"
var serviceName = "MIGRATOR"
var errDatabaseConnection = "Error connecting to db"


func main() {
	if err := run(); err != nil {
		if errors.Cause(err) != commands.ErrHelp {
			log.Printf("error: %s", err)
		}
		os.Exit(1)
	}
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
		DB   struct {
			User       string `conf:"default:rami"`
			Password   string `conf:"default:rami,noprint"`
			Host       string `conf:"default:postgres-db"` // docker-compose database service name
			Name       string `conf:"default:apimain"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	const prefix = "MIGRATETOR"
	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
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

	//
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

	if err := commands.Migrate(dbConfig); err != nil {
		return errors.Wrap(err, "migrating database")
	}

	return nil
}
