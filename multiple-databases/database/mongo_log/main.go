package mongo_log

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type WrapDB struct {
	DB          *mongo.Database
	Collections *MongoCollections
}

func initCollections(db *mongo.Database) *MongoCollections {
	collectionNames := getCollectionNames()

	logCompressionJob := db.Collection(collectionNames.LogCompressionJob)
	logCompressionTask := db.Collection(collectionNames.LogCompressionTask)
	logJobExecution := db.Collection(collectionNames.LogJobExecution)

	return &MongoCollections{
		LogCompressionJob:  logCompressionJob,
		LogCompressionTask: logCompressionTask,
		LogJobExecution:    logJobExecution,
	}
}

func InitDatabase(env *config.EnvironmentVariable) *WrapDB {
	db := newDBConn(env)
	collections := initCollections(db)

	return &WrapDB{
		DB:          db,
		Collections: collections,
	}
}
