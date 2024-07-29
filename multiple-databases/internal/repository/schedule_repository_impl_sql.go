package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"

	"github.com/rs/zerolog/log"
)

type ScheduleRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewScheduleRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) ScheduleRepository {
	return &ScheduleRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

/* Start of General Interface Impl */

func (r *ScheduleRepositoryImpl) thisDB() *sql.DB {
	return r.WrapDB.SqliteApp.DB
}

func (r *ScheduleRepositoryImpl) thisTable() string {
	return r.WrapDB.SqliteApp.Tables.CronAutoCompress
}

/* End of General Interface Impl */

func (r *ScheduleRepositoryImpl) runQuery(query string) (result model.CronAutoCompresses, err error) {
	// log.Debug().Str("query", query).Msg("runQuery")

	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	rows, err := r.thisDB().QueryContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrExecuteQuery, query)
		return
	}
	defer rows.Close()

	if rows == nil {
		return
	}

	var record model.CronAutoCompress
	for rows.Next() {
		err = rows.Scan(
			&record.Id,
			&record.CompressionId,
			&record.UserId,
			&record.EntryId,
			&record.Spec,
			&record.IsActive,
			&record.Url,
			&record.CronDetail.Format,
			&record.CronDetail.Type,
			&record.CronDetail.Repetition,
			&record.CronDetail.Value,
			&record.CreatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("rows.Scan failed")
			return
		}

		result = append(result, record)
	}

	return
}

func (r *ScheduleRepositoryImpl) Insert(req model.CronAutoCompress) error {
	// log.Info().Msg("ScheduleRepositoryImpl.Insert")

	query := fmt.Sprintf("INSERT INTO %s (compression_id, user_id, entry_id, spec, "+
		" is_active, url, cron_format, cron_type, cron_repetition, "+
		" cron_value, created_at) "+
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", r.thisTable())

	args := []any{req.CompressionId, req.UserId, req.EntryId, req.Spec,
		req.IsActive, req.Url, req.CronDetail.Format, req.CronDetail.Type, req.CronDetail.Repetition,
		req.CronDetail.Value, req.CreatedAt}

	err := ExecQuery(DBStatement{
		db:      r.thisDB(),
		timeout: r.Env.Database.Timeout.Write,
		query:   query,
		args:    args,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("statement.Exec failed")
	}

	return nil
}

func (r *ScheduleRepositoryImpl) FindByCompressionId(id string) (res model.CronAutoCompress, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE compression_id='%s' LIMIT 1 ", r.thisTable(), id)

	results, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg(lib.ErrRunQuery)
		return
	}

	if results != nil {
		res = results[0]
	}

	return
}

func (r *ScheduleRepositoryImpl) UpdateCronEntryIdByCompressionId(compressionId string, entryId int) error {
	isActive := 0
	if entryId > 0 {
		isActive = 1
	}

	args := []any{entryId, isActive, compressionId}
	query := fmt.Sprintf("UPDATE %s SET entry_id = ?, is_active = ? WHERE compression_id= ? ", r.thisTable())

	err := ExecQuery(DBStatement{
		db:      r.thisDB(),
		timeout: r.Env.Database.Timeout.Write,
		query:   query,
		args:    args,
	})

	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("execQuery failed")
		return err
	}

	return nil
}

func (r *ScheduleRepositoryImpl) FindAll() (res model.CronAutoCompresses, err error) {
	query := fmt.Sprintf("SELECT * FROM %s ", r.thisTable())

	results, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg(lib.ErrRunQuery)
		return
	}

	res = results

	return
}

func (r *ScheduleRepositoryImpl) RemoveById(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ? ", r.thisTable())

	args := []any{id}
	err := ExecQuery(DBStatement{
		db:      r.thisDB(),
		timeout: r.Env.Database.Timeout.Write,
		query:   query,
		args:    args,
	})

	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("execQuery failed")
		return err
	}

	return nil
}
