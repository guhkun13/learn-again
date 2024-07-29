package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"

	"github.com/rs/zerolog/log"
)

type LockCompressionRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewLockCompressionRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) LockCompressionRepository {
	return &LockCompressionRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

/* Start of General Interface Impl */

func (r *LockCompressionRepositoryImpl) thisDB() *sql.DB {
	return r.WrapDB.SqliteLog.DB
}

func (r *LockCompressionRepositoryImpl) thisTable() string {
	return r.WrapDB.SqliteLog.Tables.LockCompression
}

/* End of General Interface Impl */

func (r *LockCompressionRepositoryImpl) runQuery(query string) (result model.LockCompressions, err error) {
	// log.Debug().Str("query", query).Msg("runQuery")

	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	rows, err := r.thisDB().QueryContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg(lib.ErrRunQuery)
		return
	}
	defer rows.Close()

	if rows == nil {
		log.Warn().Msg("rows nil")
		return
	}

	res := model.LockCompression{}
	for rows.Next() {
		err = rows.Scan(
			&res.Id,
			&res.CompressionId,
			&res.JobId,
			&res.IsRunning,
			&res.Message,
			&res.StartDate,
			&res.EndDate,
			&res.CreatedAt,
		)

		if err != nil {
			log.Error().Err(err).Msg("rows.Scan failed")
			return
		}
		result = append(result, res)
	}

	return
}

func (r *LockCompressionRepositoryImpl) Insert(req model.LockCompression) error {
	log.Info().Msg("LockCompressionRepositoryImpl.Insert")

	query := fmt.Sprintf("INSERT INTO %s (compression_id, job_id, "+
		" is_running, message, start_date, created_at) "+
		" VALUES (?, ?, ?, ?, ?, ?)", r.thisTable())

	args := []any{req.CompressionId, req.JobId, req.IsRunning,
		req.Message, req.StartDate, req.CreatedAt}

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

func (r *LockCompressionRepositoryImpl) FindByCompressionId(id string) (res model.LockCompression, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE compression_id='%s' LIMIT 1 ", r.thisTable(), id)
	// log.Info().Msg(query)

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

func (r *LockCompressionRepositoryImpl) FindOngoingCompression(id string) (res model.LockCompression, err error) {
	// log.Info().Msg("LockCompressionRepositoryImpl.FindOngoingCompression")

	query := fmt.Sprintf("SELECT * FROM %s WHERE compression_id='%s' "+
		" AND is_running = 1 ORDER BY created_at DESC LIMIT 1 ", r.thisTable(), id)

	results, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("runQuery failed")
		return
	}

	if results != nil {
		res = results[0]
	}

	return
}

func (r *LockCompressionRepositoryImpl) FindOngoingCompressions() (results model.LockCompressions, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE is_running = 1 ORDER BY created_at DESC", r.thisTable())

	results, err = r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("runQuery failed")
		return
	}

	return
}

func (r *LockCompressionRepositoryImpl) RemoveByCompressionId(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE compression_id = ? ", r.thisTable())

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

func (r *LockCompressionRepositoryImpl) FindAll() (results model.LockCompressions, err error) {
	query := fmt.Sprintf("SELECT * FROM %s ", r.thisTable())

	results, err = r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("runQuery failed")
		return
	}

	return
}

func (r *LockCompressionRepositoryImpl) UpdateSetFinished(req model.LockCompression) error {
	// log.Info().Msg("LockCompressionRepositoryImpl.UpdateSetFinished")

	query := fmt.Sprintf("UPDATE %s SET is_running = 0, message = ?, end_date = ? "+
		" WHERE compression_id = ? AND is_running = 1 ", r.thisTable())

	args := []any{
		req.Message,
		time.Now().UTC(),
		req.CompressionId,
	}

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

func (r *LockCompressionRepositoryImpl) FindNewestRecord() (res model.LockCompression, err error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC LIMIT 1 ", r.thisTable())

	results, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("runQuery failed")
		return
	}

	if results != nil {
		res = results[0]
	}

	return
}

// update all records with running = 1 to 0, with  message="force-clear"
func (r *LockCompressionRepositoryImpl) ClearLock() error {
	query := fmt.Sprintf("UPDATE %s SET is_running = 0, message = ?, end_date = ? "+
		" WHERE is_running = 1 ", r.thisTable())

	message := "force clear"
	args := []any{message, time.Now().UTC()}

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
