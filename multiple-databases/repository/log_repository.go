package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/dto"
	"github.com/guhkun13/learn-again/multiple-databases/model"
)

type LogRepository interface {
	Create(req dto.CreateUser) error
	List() (model.Users, error)
	Read(id int) (model.User, error)
}

type LogRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewLogRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) LogRepository {
	return &LogRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}
