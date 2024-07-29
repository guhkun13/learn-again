package multipledatabases

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"

	"github.com/guhkun13/learn-again/multiple-databases/internal/handler"
)

type Handlers struct {
	UserHandler handler.UserHandler
}

func NewHandlers(env *config.EnvironmentVariable, s Services) Handlers {
	return Handlers{
		UserHandler: handler.NewUserHandler(env),
	}
}
