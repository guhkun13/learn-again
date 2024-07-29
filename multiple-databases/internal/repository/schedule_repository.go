package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type ScheduleRepository interface {
	Insert(req model.CronAutoCompress) error
	FindByCompressionId(id string) (model.CronAutoCompress, error)
	FindAll() (model.CronAutoCompresses, error)
	UpdateCronEntryIdByCompressionId(compressionId string, entryId int) error
	RemoveById(id int) error

	runQuery(query string) (result model.CronAutoCompresses, err error)
}
