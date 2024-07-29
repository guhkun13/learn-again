package repository

import (
	"fmt"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"

	"github.com/rs/zerolog/log"
)

type LogFeederRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewLogFeederRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) LogFeederRepository {
	return &LogFeederRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

func (r *LogFeederRepositoryImpl) Insert(data model.LogFeeder) error {
	table := r.WrapDB.SqliteLog.Tables.LogFeeder
	query := fmt.Sprintf("INSERT INTO %s (last_id, count_item, start_date, end_date, created_at) VALUES (?, ?, ?, ?, ?) ", table)
	// log.Info().Msg(query)

	statement, err := r.WrapDB.SqliteLog.DB.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return err
	}

	_, err = statement.Exec(data.LastId, data.CountItem, data.StartDate, data.EndDate, data.CreatedAt)
	if err != nil {
		log.Fatal().Err(err).Msg("statement.Exec failed")
	}

	return nil
}

func (r *LogFeederRepositoryImpl) FindLastRecord() (*model.LogFeeder, error) {
	res := &model.LogFeeder{}
	table := r.WrapDB.SqliteLog.Tables.LogFeeder
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC LIMIT 1 ", table)
	// log.Info().Msg(query)

	rows, err := r.WrapDB.SqliteLog.DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return res, err
	}
	defer rows.Close()

	if rows == nil {
		return res, nil
	}

	for rows.Next() {
		err = rows.Scan(
			&res.Id,
			&res.LastId,
			&res.CountItem,
			&res.StartDate,
			&res.EndDate,
			&res.CreatedAt,
		)

		if err != nil {
			log.Error().Err(err).Msg("rows.Scan failed")
			return res, err
		}
	}

	return res, nil
}
