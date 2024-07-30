package multipledatabases

import (
	"github.com/rs/zerolog/log"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		panic(err)
	}

	config.InitLogger(env)

	// Database
	wrapDB := database.InitDatabase(env)

	// App
	app := NewApp(env, wrapDB)
	app.Run()
}
