package multipledatabases

import (
	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"

	"github.com/guhkun13/learn-again/multiple-databases/internal/repository"
)

type Repositories struct {
	CompressionRepository             repository.CompressionRepository
	ScheduleRepository                repository.ScheduleRepository
	LockCompressionRepository         repository.LockCompressionRepository
	LogCompressionJobRepository       repository.LogCompressionJobRepository
	LogCompressionTaskRepository      repository.LogCompressionTaskRepository
	LogCompressionStatisticRepository repository.LogCompressionStatisticRepository
	LogJobExecutionRepository         repository.LogJobExecutionRepository
	LogFeederRepository               repository.LogFeederRepository
}

func NewRepositories(env *config.EnvironmentVariable, wrapDB *database.WrapDB) Repositories {
	return Repositories{
		CompressionRepository:             repository.NewCompressionRepositoryImpl(env, wrapDB),
		ScheduleRepository:                repository.NewScheduleRepositoryImpl(env, wrapDB),
		LockCompressionRepository:         repository.NewLockCompressionRepositoryImpl(env, wrapDB),
		LogCompressionJobRepository:       repository.NewLogCompressionJobRepositoryImpl(env, wrapDB),
		LogCompressionTaskRepository:      repository.NewLogCompressionTaskRepositoryImpl(env, wrapDB),
		LogCompressionStatisticRepository: repository.NewLogCompressionStatisticRepositoryImpl(env, wrapDB),
		LogJobExecutionRepository:         repository.NewLogJobExecutionRepositoryImpl(env, wrapDB),
		LogFeederRepository:               repository.NewLogFeederRepositoryImpl(env, wrapDB),
	}
}
