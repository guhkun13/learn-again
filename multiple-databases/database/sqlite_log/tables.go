package sqlite_log

import (
	"database/sql"
	"fmt"

	"github.com/guhkun13/learn-again/multiple-databases/database/sql_command"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/rs/zerolog/log"
)

const (
	TableLogCompressionStatistic = "log_compression_statistic"
	TableLogFeeder               = "log_feeder"
	TableLockCompression         = "lock_compression"
)

type TablesName struct {
	LogCompressionStatistic string
	LogFeeder               string
	LockCompression         string
}

func getTablesName() *TablesName {
	return &TablesName{
		LogCompressionStatistic: TableLogCompressionStatistic,
		LogFeeder:               TableLogFeeder,
		LockCompression:         TableLockCompression,
	}
}

func createTable_LogCompressionStatistic(db *sql.DB) error {
	log.Info().Msg("createTable_LogCompressionStatistic")
	query := fmt.Sprintf("%s %s ("+
		"id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "+
		"job_id VARCHAR(100) NOT NULL, "+
		"task_id VARCHAR(100) NOT NULL, "+
		"user_id VARCHAR(100) NOT NULL, "+
		"compression_id VARCHAR(100) NOT NULL, "+
		"machine_id VARCHAR(10) NOT NULL, "+
		"compressor_id VARCHAR(50) NOT NULL, "+
		"filename VARCHAR(255) NOT NULL, "+
		"format_file VARCHAR(10) NOT NULL, "+
		"original_size FLOAT NOT NULL, "+
		"compressed_size FLOAT NOT NULL, "+
		"compressed_duration FLOAT NOT NULL, "+
		"space_saving_percentage FLOAT NOT NULL, "+
		"started_at DATETIME NOT NULL, "+
		"finished_at DATETIME NOT NULL, "+
		"created_at DATETIME NOT NULL DEFAULT CURRENT_DATETIME); ",
		sql_command.GetCommand().CreateTableIfNotExists,
		TableLogCompressionStatistic)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrExecuteQuery, query)
		return err
	}

	return nil

}

func createTable_LogFeeder(db *sql.DB) error {
	log.Info().Msg("createTable_LogFeeder")
	query := fmt.Sprintf("%s %s ("+
		"id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "+
		"last_id INTEGER, "+
		"count_item INTEGER, "+
		"start_date DATETIME NOT NULL, "+
		"end_date DATETIME NOT NULL, "+
		"created_at DATETIME NOT NULL DEFAULT CURRENT_DATETIME); ",
		sql_command.GetCommand().CreateTableIfNotExists,
		TableLogFeeder)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Error().Err(err).Msgf("%s %s", lib.ErrExecuteQuery, query)
		return err
	}

	return nil

}

func createTable_LockCompression(db *sql.DB) error {
	log.Info().Msg("createTable_LockCompression")

	query := fmt.Sprintf("%s %s ("+
		"id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "+
		"compression_id VARCHAR(25), "+
		"job_id VARCHAR(50), "+
		"is_running INTEGER NOT NULL DEFAULT 0, "+
		"message VARCHAR(100), "+
		"start_date DATETIME NOT NULL, "+
		"end_date DATETIME, "+
		"created_at DATETIME NOT NULL DEFAULT CURRENT_DATETIME); ",
		sql_command.GetCommand().CreateTableIfNotExists,
		TableLockCompression)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", lib.ErrPrepareQuery, query)
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Error().Err(err).Msgf("%s %s", lib.ErrExecuteQuery, query)
		return err
	}

	return nil

}
