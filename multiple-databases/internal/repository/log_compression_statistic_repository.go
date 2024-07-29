package repository

import (
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
)

type LogCompressionStatisticRepository interface {
	Insert(req model.LogCompressionStatistic) error
	FindAll(f dto.FilterStatistic) (model.FindAllStatisticResult, error)
	FilterForDashboard(f dto.FilterStatistic) (dto.StatisticDashboards, error)
	FindOldestRecord() (model.LogCompressionStatistic, error)
	FindNewestRecord() (model.LogCompressionStatistic, error)
	Summarize() (dto.SummaryStatisticResponse, error)
}
