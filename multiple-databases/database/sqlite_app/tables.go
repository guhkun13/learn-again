package sqlite_app

import (
	"database/sql"
	"fmt"

	"github.com/guhkun13/learn-again/multiple-databases/database/sql_command"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/rs/zerolog/log"
)

const (
	TableCompression      = "compression"
	TableCronAutoCompress = "cron_auto_compress"
)

type TablesName struct {
	CronAutoCompress string
	Compression      string
}

func getTablesName() *TablesName {
	return &TablesName{
		CronAutoCompress: TableCronAutoCompress,
		Compression:      TableCompression,
	}
}

func createTable_Compression(db *sql.DB) error {
	log.Info().Msg("createTable_Compression")

	query := fmt.Sprintf("%s %s ("+
		"id VARCHAR(50) PRIMARY KEY, "+
		"user_id INTEGER, "+
		"source_path TEXT NOT NULL, "+
		"destination_path TEXT NOT NULL, "+
		"repetition VARCHAR(50) NOT NULL, "+
		"created_at DATETIME NOT NULL DEFAULT CURRENT_DATETIME); ",
		sql_command.GetCommand().CreateTableIfNotExists,
		TableCompression)

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

func createTable_CronAutoCompress(db *sql.DB) error {
	log.Info().Msg("createTable_CronAutoCompress")

	query := fmt.Sprintf("%s %s ("+
		"id INTEGER PRIMARY KEY AUTOINCREMENT, "+
		"compression_id VARCHAR(50) NOT NULL, "+
		"user_id INTEGER, "+
		"entry_id INTEGER NOT NULL, "+
		"spec VARCHAR(50) NOT NULL, "+
		"is_active INTEGER NOT NULL, "+
		"url TEXT NOT NULL, "+
		"cron_format VARCHAR(25) NOT NULL, "+
		"cron_type VARCHAR(25) NOT NULL, "+
		"cron_repetition VARCHAR(25) NOT NULL, "+
		"cron_value VARCHAR(25), "+
		"created_at DATETIME NOT NULL DEFAULT CURRENT_DATETIME, "+
		"CONSTRAINT fk_compression_id FOREIGN KEY (compression_id) REFERENCES compression(id) ON DELETE CASCADE); ",
		sql_command.GetCommand().CreateTableIfNotExists,
		TableCronAutoCompress)

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
