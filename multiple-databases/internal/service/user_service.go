package service

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type UserService interface {
	FindAll(filter dto.FilterFindCompression) (model.FindCompressionWithPagination, error)
	FindById(id string) (model.Compression, error)
	DeleteById(id string) error

	UploadCompress(req dto.UploadCompressRequest) (string, error)
	RegisterAutoCompress(req dto.RegisterAutoCompressRequest) (model.Compression, error)
	UpdateAutoCompress(compression model.Compression, req dto.UpdateAutoCompressRequest) error

	// multiple ids
	FindByIds(ids string) (model.Compressions, error)
	DeleteByIds(ids []string) (int, error)
}
