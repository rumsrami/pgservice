package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate migrates the schema
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add title",
		Script: `
CREATE TABLE "info"
(
	"title"         text NOT NULL,
	"uuid4"         text UNIQUE NOT NULL,
	"timestamp"     timestamp without time zone NOT NULL
);`,
	},
}
