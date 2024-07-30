package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/dto"
	"github.com/guhkun13/learn-again/multiple-databases/model"
)

type ArticleRepository interface {
	Create(req dto.CreateUser) error
	List() (model.Users, error)
	Read(id int) (model.User, error)
}

type ArticleRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewArticleRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) ArticleRepository {
	return &ArticleRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}
