package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/dto"
	"github.com/guhkun13/learn-again/multiple-databases/model"
)

type UserRepository interface {
	Create(req dto.CreateUser) error
	List() (model.Users, error)
	Read(id int) (model.User, error)
}

type UserRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewUserRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) UserRepository {
	return &UserRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}
