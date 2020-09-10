package db

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // database driver
)

const (
	errDatabaseConnectionString = "cannot connect to database using this connection string"
	errDatabaseSessionFailed    = "cannot create database session"
)

// Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

// OpenTCP connects to the database
func OpenTCP(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	dbURL := make(url.Values)
	dbURL.Set("sslmode", sslMode)
	dbURL.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: dbURL.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// Ping checks if the database is on
func Ping(ctx context.Context, db *sqlx.DB) error {
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
