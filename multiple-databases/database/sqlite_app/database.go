package sqlite_app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	config "github.com/guhkun13/learn-again/multiple-databases/config"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rs/zerolog/log"
)

const DBDriver = "sqlite3"

func newDBConn(env *config.EnvironmentVariable) *sql.DB {
	dbDir := env.Database.SqliteApp.Dir
	if dbDir == "" {
		dbDir = lib.DefaultDbDirName
	}

	dbName := env.Database.SqliteApp.Name
	dbPath := fmt.Sprintf("%s/%s", dbDir, dbName)

	dsn := fmt.Sprintf("file:%s", dbPath)
	createDatabaseIfNotExist(dbPath)

	db, err := sql.Open(DBDriver, dsn)
	if err != nil {
		panic(err)
	}
	// db.SetMaxOpenConns(1)
	// defer db.Close()

	log.Info().Msgf("== SQLITE APP DB [%s] connected ==", DBDriver)

	return db
}

func createDatabaseIfNotExist(path string) {
	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msgf("Create path (%s) failed", path)
			panic(err)
		}
		file.Close()
	}
	db, err := sql.Open(DBDriver, path)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 2)

	err = createTables(db)
	if err != nil {
		log.Error().Err(err).Msg("createTables failed")
	}
}

func createTables(db *sql.DB) error {

	err := createTable_Compression(db)
	if err != nil {
		log.Error().Err(err).Msg("createTable_Compression failed")
		panic(err)
	}

	err = createTable_CronAutoCompress(db)
	if err != nil {
		log.Error().Err(err).Msg("createTable_LogCompressionStatistic failed")
		panic(err)
	}

	return nil
}
