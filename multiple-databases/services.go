package multipledatabases

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/internal/service"
)

type Services struct {
	CompressionService     service.CompressionService
	ScheduleService        service.ScheduleService
	LogJobExecutionService service.LogJobExecutionService
	LockCompressionService service.LockCompressionService
	MonitoringService      service.MonitoringService
	SummaryService         service.SummaryService
	LogFeederService       service.LogFeederService
}

func NewServices(
	env *config.EnvironmentVariable,
	r Repositories,
) Services {

	return Services{
		UserService:     service.NewCompressionServiceImpl(env, r.CompressionRepository)
	}
}
