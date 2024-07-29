package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type LogJobExecutionRepository interface {
	Insert(req model.LogJobExecution) error
	CountJobSummaryWithFilter(f dto.LogCountFilter) (int64, error)
	CountSummaryJobSteps() (model.JobSummaries, error)
}
