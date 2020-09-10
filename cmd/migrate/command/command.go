package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/rumsrami/pgservice/internal/data/schema"
	"github.com/rumsrami/pgservice/internal/db"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

// Migrate creates the schema in the database.
func Migrate(cfg db.Config) error {
	db, err := db.OpenTCP(cfg)
	if err != nil {
		return errors.Wrap(err, "connect database")
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return errors.Wrap(err, "migrate database")
	}

	fmt.Println("migrations complete")
	return nil
}