package multipledatabases

import (
	"github.com/gin-gonic/gin"

	config "github.com/guhkun13/learn-again/multiple-databases/config"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/guhkun13/learn-again/multiple-databases/internal/handler"
)

type Handler struct {
	Env         *config.EnvironmentVariable
	UserHandler handler.UserHandler
}

func NewEngine(handler Handler) *gin.Engine {
	router := gin.Default()

	// set gin mode
	if handler.Env.App.Mode == lib.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// api
	apiRouterGroup := router.Group("/api")
	rg := apiRouterGroup.Group("/v1")

	r := rg.Group("/auth")

	r.POST("/login", handler.UserHandler.Create)
	r.GET("/get-user/:username", handler.UserHandler.GetUser)

	return router
}
