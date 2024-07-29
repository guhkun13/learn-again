package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/guhkun13/learn-again/multiple-databases/internal/model"

	"github.com/rs/zerolog/log"
)

func AddWhereOrAnd(query string, isFirst bool) (string, bool) {
	if isFirst {
		query += " WHERE "
		return query, false
	}
	query += " AND "
	return query, false
}

func CountRecordFromQuery(db *sql.DB, query string) (res model.CountAll, err error) {
	query = strings.ReplaceAll(query, "SELECT * ", "SELECT COUNT(*) ")

	rows, err := db.Query(query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msgf(lib.ErrRunQuery)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			log.Error().Err(err).Msg("rows.Scan failed")
			return
		}
	}

	return
}

type DBStatement struct {
	db      *sql.DB
	timeout time.Duration
	query   string
	args    []any
}

func ExecQuery(d DBStatement) error {
	// log.Debug().Str("query", d.query).Any("args", d.args).Msg("ExecQuery")

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	statement, err := d.db.PrepareContext(ctx, d.query)
	if err != nil {
		log.Error().Err(err).Str("query", d.query).Msg(lib.ErrPrepareQuery)
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, d.args...)
	if err != nil {
		log.Error().Err(err).Str("query", d.query).Msg(lib.ErrExecuteQuery)
		return err
	}

	return nil
}
