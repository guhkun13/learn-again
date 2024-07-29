package mongo_log

import (
	"context"
	"fmt"

	"github.com/guhkun13/learn-again/multiple-databases/config"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/guhkun13/encryptor"
)

func getDbUri(env *config.EnvironmentVariable) string {
	var err error
	host := env.Database.MongoLog.Host
	port := env.Database.MongoLog.Port
	username := env.Database.MongoLog.Username
	password := env.Database.MongoLog.Password

	if env.Secret.IsEncrypted {
		password, err = encryptor.DecryptByKeyCombination(password)
		if err != nil {
			log.Error().Err(err).Msg("failed DecryptByKeyCombination")
		}
	}

	return fmt.Sprintf("mongodb://%s:%s@%s:%d/", username, password, host, port)
}

func newDBConn(env *config.EnvironmentVariable) *mongo.Database {
	dbUri := getDbUri(env)

	// define db name
	dbName := env.Database.MongoLog.Name

	// init context
	ctx, cancel := context.WithTimeout(context.Background(), env.Database.Timeout.Ping)
	defer cancel()

	// Connect to the database.
	clientOptions := options.Client().ApplyURI(dbUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed connect to mongodb [%s]", dbName)
	}

	// Check the connection.
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("Not Connected to mongodb [%s]", dbName)
	}
	log.Info().Msgf("== mongoDB [%s] connected ==", dbName)

	// establish connection
	return client.Database(dbName)
}
