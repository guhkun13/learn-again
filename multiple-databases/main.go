package multipledatabases

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/internal/handler"
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

type Handler struct {
	Env         *config.EnvironmentVariable
	UserHandler handler.UserHandler
}

func (s *AppImpl) Run() {
	var err error

	// Repo
	repositories := NewRepositories(s.Env, s.WrapDb)

	// Service
	services := NewServices(s.Env, repositories)

	// Handler
	handlers := NewHandlers(s.Env, services)

	routes := Handler{
		Env:         s.Env,
		UserHandler: handlers.UserHandler,
	}

	// define routes
	addr := fmt.Sprintf("%s:%d", s.Env.App.Host, s.Env.App.Port)

	router := gin.Default()

	// api
	apiRouterGroup := router.Group("/api")
	rg := apiRouterGroup.Group("/v1")

	r := rg.Group("/auth")

	r.POST("/login", routes.UserHandler.Create())
	r.GET("/get-user/:username", handler.UserHandler.GetUser)

	log.Info().Str("address", addr).Msg("Run App")

	go func() {
		err = router.Run(addr)
		if err != nil {
			log.Fatal().Err(err).Msg("XXX [Fatal Error] XXX")
		}
	}()

}
