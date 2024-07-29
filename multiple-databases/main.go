package multipledatabases

import (
	"fmt"

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

type AppImpl struct {
	Env    *config.EnvironmentVariable
	WrapDb *database.WrapDB
}

func NewApp(env *config.EnvironmentVariable,
	wrapDb *database.WrapDB,
) *AppImpl {
	return &AppImpl{
		Env:    env,
		WrapDb: wrapDb,
	}
}

func (s *AppImpl) Run() {
	var err error

	// Repo
	repositories := NewRepositories(s.Env, s.WrapDb)

	// Service
	services := NewServices(s.Env, repositories)

	handlers := NewHandlers(s.Env, services)

	routes := Handler{
		Env:         s.Env,
		UserHandler: handlers.UserHandler,
	}

	// define routes
	addr := fmt.Sprintf("%s:%d", s.Env.App.Host, s.Env.App.Port)
	engine := NewEngine(routes)

	log.Info().Str("address", addr).Msg("Run App")

	// start server as go-routine as to not block next code execution: running nsq consumers
	go func() {
		err = engine.Run(addr)
		if err != nil {
			log.Fatal().Err(err).Msg("XXX [Fatal Error] XXX")
		}
	}()

}
