package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type LogCompressionJobRepository interface {
	Insert(req model.CompressionJob) error
	FindAll(filter dto.Filtering) (model.FindAllJobResult, error)
	FindById(id string) (model.CompressionJob, error)
	FindByCompressionId(id string) (model.CompressionJob, error)
}
