package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/guhkun13/learn-again/multiple-databases/dto"
	"github.com/guhkun13/learn-again/multiple-databases/model"

	"github.com/rs/zerolog/log"
)

func (r *UserRepositoryImpl) runQuery(query string) (model.Users, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	rows, err := r.WrapDB.SqliteApp.DB.QueryContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to read")
		return nil, err
	}
	defer rows.Close()

	if rows == nil {
		log.Warn().Msg("rows nil")
		return nil, nil
	}

	var result model.Users
	res := model.User{}
	for rows.Next() {
		err = rows.Scan(
			&res.ID,
			&res.Name,
			&res.Email,
		)

		if err != nil {
			log.Error().Err(err).Msg("rows.Scan failed")
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

func (r *UserRepositoryImpl) Create(data model.User) error {
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

func (r *UserRepositoryImpl) Read(id int) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s ", r.thisTable())

	result, err = r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("runQuery failed")
		return
	}

	return
}

func (r *UserRepositoryImpl) List() (model.Users, error) {
	query := fmt.Sprintf("SELECT * FROM %s ", r.thisTable())

	result, err = r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("runQuery failed")
		return
	}

	return
}

func (r *UserRepositoryImpl) Update(data model.User, req dto.UpdateUser) error {
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

func (r *UserRepositoryImpl) Delete(data model.User) error {
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
