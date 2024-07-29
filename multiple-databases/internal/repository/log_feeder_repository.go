package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type LogFeederRepository interface {
	Insert(req model.LogFeeder) error
	FindLastRecord() (*model.LogFeeder, error)
}
