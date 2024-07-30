package multipledatabases

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/internal/service"
)

type Services struct {
	UserService     service.UserService
}

func NewServices(
	env *config.EnvironmentVariable,
	r Repositories,
) Services {
	return Services{
		UserService:     service.NewCompressionServiceImpl(env, r.CompressionRepository)
	}
}
