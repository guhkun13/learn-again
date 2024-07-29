package sqlite_log

import (
	"database/sql"

	"github.com/guhkun13/learn-again/multiple-databases/config"
)

type WrapDB struct {
	DB     *sql.DB
	Tables *TablesName
}

func InitDatabase(env *config.EnvironmentVariable) *WrapDB {
	db := newDBConn(env)
	tablesName := getTablesName()

	return &WrapDB{
		DB:     db,
		Tables: tablesName,
	}
}
