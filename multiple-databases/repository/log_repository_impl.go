package repository

import (
	"context"
	"fmt"

	"github.com/guhkun13/learn-again/multiple-databases/dto"
	"github.com/guhkun13/learn-again/multiple-databases/model"

	"github.com/rs/zerolog/log"
)

func (r *LogRepositoryImpl) runQuery(query string) (model.Users, error) {
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
			&res.Age,
		)

		if err != nil {
			log.Error().Err(err).Msg("rows.Scan failed")
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

func (r *LogRepositoryImpl) Create(req dto.CreateUser) error {
	query := fmt.Sprintf("INSERT INTO %s (name, age) VALUES (?, ?)", r.WrapDB.SqliteApp.Tables.User)

	args := []any{req.Name, req.Age}

	err := ExecQuery(DBStatement{
		db:      r.WrapDB.SqliteApp.DB,
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

func (r *LogRepositoryImpl) List() (model.Users, error) {
	query := fmt.Sprintf("SELECT * FROM %s ", r.WrapDB.SqliteApp.Tables.User)

	res, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("runQuery failed")
		return nil, err
	}

	return res, nil
}

func (r *LogRepositoryImpl) Read(id int) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", r.WrapDB.SqliteApp.Tables.User, id)

	res, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("runQuery failed")
		return model.User{}, err
	}

	return res[0], nil
}
