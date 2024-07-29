package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type LogCompressionTaskRepository interface {
	Insert(req model.LogCompressionTask) error
	FindById(id string) (model.LogCompressionTask, error)
	FindAll(f dto.Filtering) (model.FindAllTaskResult, error)
	FindAllByJobId(jobId string, filter dto.Filtering) (model.FindAllTaskResult, error)
	CountTaskSummary() (res model.TaskSummaries, err error)
}
