package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type create struct{}

// Create ...
var Create create

func (create) Info(ctx context.Context, db *sqlx.DB, title string) (*Info, error) {
	i := Info{
		Title:     title,
		UUID4:    uuid.New().String(),
		Timestamp: TimeNowUTC(),
	}

	const q = `INSERT INTO info
	(title, uuid4, timestamp)
	VALUES ($1, $2, $3)`
	_, err := db.ExecContext(ctx, q, i.Title, i.UUID4, i.Timestamp)
	if err != nil {
		return nil, errors.Wrap(err, "inserting user")
	}

	return &i, nil
}
