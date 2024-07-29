package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type LockCompressionRepository interface {
	Insert(req model.LockCompression) error
	FindByCompressionId(compressionId string) (model.LockCompression, error)
	FindOngoingCompression(compressionId string) (model.LockCompression, error)
	FindOngoingCompressions() (model.LockCompressions, error)
	FindAll() (model.LockCompressions, error)

	UpdateSetFinished(req model.LockCompression) error
	FindNewestRecord() (model.LockCompression, error)

	ClearLock() error

	runQuery(query string) (result model.LockCompressions, err error)
}
