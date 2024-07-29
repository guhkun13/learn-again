package repository

import (
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"

	config "github.com/guhkun13/learn-again/multiple-databases/config"
	database "github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/database/sqlite_log"
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"
)

type LogCompressionStatisticRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewLogCompressionStatisticRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) LogCompressionStatisticRepository {
	return &LogCompressionStatisticRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

func (r *LogCompressionStatisticRepositoryImpl) ThisDB() *sqlite_log.WrapDB {
	return r.WrapDB.SqliteLog
}

func (r *LogCompressionStatisticRepositoryImpl) ThisTable() string {
	return r.ThisDB().Tables.LogCompressionStatistic
}

func (r *LogCompressionStatisticRepositoryImpl) executeQuery(query string) (result model.LogCompressionStatistics, err error) {
	rows, err := r.ThisDB().DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrExecuteQuery, query)
		return
	}
	defer rows.Close()

	if rows == nil {
		return
	}

	var record model.LogCompressionStatistic
	for rows.Next() {
		err = rows.Scan(
			&record.Id,
			&record.JobId,
			&record.TaskId,
			&record.UserId,
			&record.CompressionId,
			&record.MachineId,
			&record.CompressorId,
			&record.Filename,
			&record.FormatFile,
			&record.OriginalSize,
			&record.CompressedSize,
			&record.CompressedDuration,
			&record.SpaceSavingPercentage,
			&record.StartedAt,
			&record.FinishedAt,
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

func (r *LogCompressionStatisticRepositoryImpl) Insert(req model.LogCompressionStatistic) error {
	// log.Info().Msg("!!! LogCompressionStatisticRepositoryImpl.Insert !!!")

	query := fmt.Sprintf("INSERT INTO %s (job_id, task_id, user_id, "+
		" compression_id, machine_id, compressor_id, filename,"+
		" format_file, original_size, compressed_size, compressed_duration,"+
		" space_saving_percentage, started_at, finished_at, created_at) "+
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.ThisTable())

	// log.Info().
	// Interface("req", req).
	// Msg(query)

	statement, err := r.ThisDB().DB.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return err
	}

	_, err = statement.Exec(req.JobId, req.TaskId, req.UserId,
		req.CompressionId, req.MachineId, req.CompressorId, req.Filename,
		req.FormatFile, req.OriginalSize, req.CompressedSize, req.CompressedDuration,
		req.SpaceSavingPercentage, req.StartedAt, req.FinishedAt, req.CreatedAt)
	if err != nil {
		log.Fatal().Err(err).Msg("statement.Exec failed")
	}

	return nil
}

func (r *LogCompressionStatisticRepositoryImpl) CountAll() (res model.CountAll, err error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s ", r.ThisTable())

	rows, err := r.ThisDB().DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
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

func (r *LogCompressionStatisticRepositoryImpl) FindAll(f dto.FilterStatistic) (res model.FindAllStatisticResult, err error) {
	// log.Debug().Interface("f", f).Msg("LogCompressionStatisticRepositoryImpl.FindAll")
	var records model.LogCompressionStatistics

	sortBy := "finished_at"
	sortOrder := "DESC"

	if f.Filtering.Sort.SortOrder != "" {
		sortOrder = f.Filtering.Sort.SortOrder
	}
	if f.Filtering.Sort.SortBy != "" {
		sortBy = f.Filtering.Sort.SortBy
	}

	size := f.Filtering.Pagination.Limit
	page := f.Filtering.Pagination.Page

	// log.Info().
	// 	Interface("timestamp", f.Timestamp).
	// 	Interface("startAt", startAt).
	// 	Interface("endAt", endAt).
	// 	Msg("filter")

	query := fmt.Sprintf("SELECT * FROM %s ", r.ThisTable())

	isFirst := true
	if f.Timestamp != nil {
		startAt := f.Timestamp.Start.Format(time.RFC3339)
		endAt := f.Timestamp.End.Format(time.RFC3339)
		if startAt != "" && endAt != "" {
			query, isFirst = AddWhereOrAnd(query, isFirst)
			query += fmt.Sprintf(" finished_at >= '%s' AND finished_at < '%s' ", startAt, endAt)
		}
	}

	if f.Column.JobId != "" {
		query, isFirst = AddWhereOrAnd(query, isFirst)
		query += fmt.Sprintf(" job_id='%s' ", f.Column.JobId)
	}

	if f.Column.CompressionId != "" {
		query, _ = AddWhereOrAnd(query, isFirst)
		query += fmt.Sprintf(" compression_id='%s' ", f.Column.CompressionId)
	}

	total, err := CountRecordFromQuery(r.ThisDB().DB, query)
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

	records, err = r.executeQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("executeQuey failed")
		return
	}

	totalPage := 1
	if size > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(size)))
	}

	res.LogCompressionStatistics = records
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

func (r *LogCompressionStatisticRepositoryImpl) FilterForDashboard(f dto.FilterStatistic) (res dto.StatisticDashboards, err error) {
	// log.Info().
	// 	Interface("filter", f).
	// 	Msg("LogCompressionStatisticRepositoryImpl.FilterForDashboard")

	/*
		SELECT
		SUM(original_size) AS total_original_size,
		SUM(compressed_size) AS total_compressed_size,
		strftime('%Y-%m-%d %H',finished_at) ts
		FROM log_compression_statistic lcs
		WHERE finished_at >= "2023-04-22"
			AND finished_at <= "2025-04-24"
		GROUP BY strftime('%H',finished_at)
		ORDER BY finished_at DESC
		LIMIT 24
	*/

	displayDate := "%Y-%m-%d"
	groupBy := "%d"
	limit := 7
	if f.GroupBy.Value == lib.GroupByHour {
		displayDate += " %H"
		groupBy = "%H"
		limit = 24
	}

	query := fmt.Sprintf("SELECT"+
		" SUM(original_size) AS total_original_size, "+
		" SUM(compressed_size) AS total_compressed_size, "+
		" strftime('%s',finished_at) ts "+
		" FROM %s "+
		" WHERE finished_at >= '%s' "+
		" AND finished_at <= '%s' "+
		" GROUP BY strftime('%s',finished_at) "+
		" ORDER BY finished_at DESC"+
		" LIMIT %d",
		displayDate,
		r.ThisTable(),
		f.Timestamp.Start.Format(time.RFC3339),
		f.Timestamp.End.Format(time.RFC3339),
		groupBy, limit)

	// log.Info().Str("query", query).Msg("query")
	rows, err := r.ThisDB().DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return
	}
	defer rows.Close()

	var record dto.StatisticDashboard
	for rows.Next() {
		err = rows.Scan(
			&record.TotalOriginalSize,
			&record.TotalCompressedSize,
			&record.Timestamp,
		)

		if err != nil {
			log.Error().Err(err).Msg("response.Scan failed")
			return
		}

		res = append(res, record)
	}

	return res, nil
}

func (r *LogCompressionStatisticRepositoryImpl) Summarize() (res dto.SummaryStatisticResponse, err error) {
	// log.Info().Msg("LogCompressionStatisticRepositoryImpl.Summarize")

	/*
		total_files:
		total_compressed_duration:
		total_original_size:
		total_compressed_size:

		SELECT COUNT(*) as total_files,
		SUM(original_size) as total_original_size,
		SUM(compressed_size) as total_compressed_size,
		SUM(compressed_duration) as total_compressed_duration
		(1 - (SUM(compressed_size)/SUM(original_size))) * 100 as total_space_saving_percentage
		FROM log_compression_statistic lcs
	*/

	query := fmt.Sprintf("SELECT COUNT(*) as total_files,"+
		" SUM(original_size) as total_original_size,"+
		" SUM(compressed_size) as total_compressed_size,"+
		" SUM(compressed_duration) as total_compressed_duration,"+
		" (1 - (SUM(compressed_size)/SUM(original_size))) * 100 as total_space_saving_percentage "+
		" FROM %s",
		r.ThisTable())

	rows, err := r.ThisDB().DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&res.TotalFiles,
			&res.TotalOriginalSize,
			&res.TotalCompressedSize,
			&res.TotalCompressedDuration,
			&res.TotalSpaceSavingPercentage)

		if err != nil {
			log.Error().Err(err).Msg("response.Scan failed")
			return
		}
	}

	return
}

func (r *LogCompressionStatisticRepositoryImpl) FindOldestRecord() (res model.LogCompressionStatistic, err error) {
	table := r.ThisTable()
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at ASC LIMIT 1", table)
	log.Info().Msg(query)

	results, err := r.executeQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("executeQuey failed")
		return
	}

	if results != nil {
		res = results[0]
	}

	return
}

func (r *LogCompressionStatisticRepositoryImpl) FindNewestRecord() (res model.LogCompressionStatistic, err error) {
	table := r.ThisTable()
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC LIMIT 1 ", table)
	log.Info().Msg(query)

	results, err := r.executeQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("executeQuey failed")
		return
	}

	if results != nil {
		res = results[0]
	}

	return
}
