package database

import (
	config "github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database/mongo_log"
	"github.com/guhkun13/learn-again/multiple-databases/database/sqlite_app"
	"github.com/guhkun13/learn-again/multiple-databases/database/sqlite_log"
)

type WrapDB struct {
	// MongoApp  *mongo_app.WrapDB
	MongoLog  *mongo_log.WrapDB
	SqliteApp *sqlite_app.WrapDB
	SqliteLog *sqlite_log.WrapDB
}

func InitDatabase(env *config.EnvironmentVariable) *WrapDB {
	// mongoAppDB := mongo_app.InitDatabase(env)
	mongoLogDB := mongo_log.InitDatabase(env)
	sqliteAppDB := sqlite_app.InitDatabase(env)
	sqliteLogDB := sqlite_log.InitDatabase(env)

	return &WrapDB{
		// MongoApp:  mongoAppDB,
		MongoLog:  mongoLogDB,
		SqliteApp: sqliteAppDB,
		SqliteLog: sqliteLogDB,
	}
}
