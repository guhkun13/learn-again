package multipledatabases

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/repository"
)

type Repositories struct {
	User    repository.UserRepository
	Article repository.ArticleRepository
	Log     repository.LogRepository
}

func NewRepositories(env *config.EnvironmentVariable, wrapDB *database.WrapDB) *Repositories {
	return &Repositories{
		User:    repository.NewUserRepositoryImpl(env, wrapDB),
		Article: repository.NewArticleRepositoryImpl(env, wrapDB),
		Log:     repository.NewLogRepositoryImpl(env, wrapDB),
	}
}
