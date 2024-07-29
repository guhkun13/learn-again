package mongo_log

import "go.mongodb.org/mongo-driver/mongo"

type CollectionNames struct {
	LogCompressionJob  string
	LogCompressionTask string
	LogJobExecution    string
}

func getCollectionNames() CollectionNames {
	return CollectionNames{
		LogCompressionJob:  "log_compression_job",
		LogCompressionTask: "log_compression_task",
		LogJobExecution:    "log_job_execution",
	}
}

type MongoCollections struct {
	LogCompressionJob  *mongo.Collection
	LogCompressionTask *mongo.Collection
	LogJobExecution    *mongo.Collection
}
