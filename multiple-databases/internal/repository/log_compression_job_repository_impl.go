package repository

import (
	"context"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/database"

	"github.com/guhkun13/learn-again/multiple-databases/config"
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/helper"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type LogCompressionJobRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewLogCompressionJobRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,

) LogCompressionJobRepository {
	return &LogCompressionJobRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

func (r *LogCompressionJobRepositoryImpl) Insert(req model.CompressionJob) error {

	_, err := r.WrapDB.MongoLog.Collections.LogCompressionJob.InsertOne(context.TODO(), req)
	if err != nil {
		return err
	}
	return nil
}

func (r *LogCompressionJobRepositoryImpl) FindAll(f dto.Filtering) (res model.FindAllJobResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	filter := bson.M{}

	if f.Search.CompressionId != "" {
		filter["compression_id"] = f.Search.CompressionId
	}

	log.Debug().Interface("filter", f).Msg("LogCompressionJobRepositoryImpl.FindAll")

	var compressionJobs model.CompressionJobs

	pgData, err := mongopagination.New(r.WrapDB.MongoLog.Collections.LogCompressionJob).Context(ctx).
		Limit(int64(f.Pagination.Limit)).
		Page(int64(f.Pagination.Page)).
		Sort(f.Sort.SortBy, helper.ConvertSortOrderDirectionForMongo(f.Sort.SortOrder)).
		Filter(filter).
		Decode(&compressionJobs).Find()

	if err != nil {
		log.Error().Err(err).Msg("err on pgdata")
		return
	}

	res.CompressionJobs = compressionJobs
	res.PaginatedData = pgData

	return
}

func (r *LogCompressionJobRepositoryImpl) FindById(id string) (res model.CompressionJob, err error) {
	filter := bson.M{"id": id}

	err = r.WrapDB.MongoLog.Collections.LogCompressionJob.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return
	}

	return
}

func (r *LogCompressionJobRepositoryImpl) FindByCompressionId(id string) (res model.CompressionJob, err error) {
	filter := bson.M{"compression_id": id}

	err = r.WrapDB.MongoLog.Collections.LogCompressionJob.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return
	}

	return
}
