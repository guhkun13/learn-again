package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	config "github.com/guhkun13/learn-again/multiple-databases/config"
	database "github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/rs/zerolog/log"
)

type CompressionRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewCompressionRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) CompressionRepository {
	return &CompressionRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

/* Start of General Interface Impl */

func (r *CompressionRepositoryImpl) thisDB() *sql.DB {
	return r.WrapDB.SqliteApp.DB
}

func (r *CompressionRepositoryImpl) thisTable() string {
	return r.WrapDB.SqliteApp.Tables.Compression
}

/* End of General Interface Impl */

func (r *CompressionRepositoryImpl) runQuery(query string) (result model.Compressions, err error) {
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

	var record model.Compression
	for rows.Next() {
		err = rows.Scan(
			&record.Id,
			&record.UserId,
			&record.SourcePath,
			&record.DestinationPath,
			&record.Repetition,
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

func (r *CompressionRepositoryImpl) Insert(req model.Compression) error {
	// log.Info().Msg("CompressionRepositoryImpl.Insert")

	query := fmt.Sprintf("INSERT INTO %s (id, user_id, source_path, destination_path, repetition, created_at)"+
		" VALUES (?, ?, ?, ?, ?, ?)", r.thisTable())

	args := []any{req.Id, req.UserId, req.SourcePath, req.DestinationPath, string(req.Repetition), req.CreatedAt}

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

func (r *CompressionRepositoryImpl) FindAll(f dto.FilterFindCompression) (res model.FindCompressionWithPagination, err error) {
	// log.Debug().Interface("f", f).Msg("CompressionRepositoryImpl.FindAll")
	var records model.Compressions

	sortBy := "created_at"
	sortOrder := "DESC"

	if f.Filtering.Sort.SortOrder != "" {
		sortOrder = f.Filtering.Sort.SortOrder
	}

	if f.Filtering.Sort.SortBy != "" {
		sortBy = f.Filtering.Sort.SortBy
	}

	size := f.Filtering.Pagination.Limit
	page := f.Filtering.Pagination.Page

	query := fmt.Sprintf("SELECT * FROM %s ", r.thisTable())

	isFirst := true

	if f.Repetition != "" {
		query, isFirst = AddWhereOrAnd(query, isFirst)
		query += fmt.Sprintf(" repetition='%s' ", f.Repetition)
	}

	total, err := CountRecordFromQuery(r.thisDB(), query)
	if err != nil {
		return
	}

	if f.Filtering.Sort.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY LOWER(%s) %s ", sortBy, sortOrder)
	}

	if size > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d",
			size, (page-1)*size)
	}

	records, err = r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("executeQuey failed")
		return
	}

	totalPage := 1
	if size > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(size)))
	}

	res.Compressions = records
	res.DBPagination = model.DBPagination{
		Total:     int64(total),
		Page:      int64(f.Filtering.Pagination.Page),
		Size:      int64(f.Filtering.Pagination.Limit),
		Prev:      int64(f.Filtering.Pagination.Page) - 1,
		Next:      int64(f.Filtering.Pagination.Page) + 1,
		TotalPage: int64(totalPage),
	}
	return
}

func (r *CompressionRepositoryImpl) FindById(id string) (res model.Compression, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id='%s'", r.thisTable(), id)

	results, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("executeQuey failed")
		return
	}

	if results != nil {
		res = results[0]
	}

	return

}

func (r *CompressionRepositoryImpl) UpdateAutoCompressPath(compression model.Compression, req dto.UpdateAutoCompressRequest) error {
	args := []any{req.SourcePath, req.DestinationPath, compression.Id}
	query := fmt.Sprintf("UPDATE %s SET source_path = ?, destination_path = ? WHERE id=? ", r.thisTable())

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

func (r *CompressionRepositoryImpl) RemoveById(id string) error {
	args := []any{id}
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (?) ", r.thisTable())

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

// Multiple

func (r *CompressionRepositoryImpl) FindByIds(ids []string) (res model.Compressions, err error) {
	strIds := lib.ListToStringWithSingleQuote(ids)
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", r.thisTable(), strIds)

	results, err := r.runQuery(query)
	if err != nil {
		log.Error().Err(err).Msg(lib.ErrRunQuery)
		return
	}

	res = results

	return
}

func (r *CompressionRepositoryImpl) RemoveByIds(ids []string) (countDeleted int, err error) {
	strIds := lib.ListToStringWithSingleQuote(ids)
	args := []any{}

	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (%s) ", r.thisTable(), strIds)

	err = ExecQuery(DBStatement{
		db:      r.thisDB(),
		timeout: r.Env.Database.Timeout.Write,
		query:   query,
		args:    args,
	})

	if err != nil {
		log.Error().Err(err).Str("query", query).Msg(lib.ErrExecuteQuery)
		return
	}

	return
}
