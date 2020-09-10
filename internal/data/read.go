package data

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type read struct{}

// Read ...
var Read read

func (read) Info(ctx context.Context, db *sqlx.DB, title string) (*Info, error) {
	var i Info
	
	const q = `SELECT * FROM info WHERE title = $1`
	if err := db.GetContext(ctx, &i, q, title); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrapf(err, "selecting info %q", title)
	}

	return &i, nil
}
