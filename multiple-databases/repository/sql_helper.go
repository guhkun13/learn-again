package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
)

type DBStatement struct {
	db      *sql.DB
	timeout time.Duration
	query   string
	args    []any
}

func AddWhereOrAnd(query string, isFirst bool) (string, bool) {
	if isFirst {
		query += " WHERE "
		return query, false
	}
	query += " AND "
	return query, false
}

func ExecQuery(d DBStatement) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	statement, err := d.db.PrepareContext(ctx, d.query)
	if err != nil {
		log.Error().Err(err).Str("query", d.query).Msg("Failed to prepare query")
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, d.args...)
	if err != nil {
		log.Error().Err(err).Str("query", d.query).Msg("Failed to execute query")
		return err
	}

	return nil
}
