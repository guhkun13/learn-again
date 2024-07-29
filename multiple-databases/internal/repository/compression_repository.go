package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type CompressionRepository interface {
	Insert(req model.Compression) error
	UpdateAutoCompressPath(compression model.Compression, req dto.UpdateAutoCompressRequest) error
	FindById(id string) (model.Compression, error)
	RemoveById(id string) error
	FindAll(filter dto.FilterFindCompression) (model.FindCompressionWithPagination, error)
	FindByIds(id []string) (model.Compressions, error)
	RemoveByIds(id []string) (int, error)

	runQuery(query string) (result model.Compressions, err error)
}
