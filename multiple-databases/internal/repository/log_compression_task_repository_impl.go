package repository

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	config "github.com/guhkun13/learn-again/multiple-databases/config"
	database "github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/helper"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"
)

type LogCompressionTaskRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewLogCompressionTaskRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) LogCompressionTaskRepository {
	return &LogCompressionTaskRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

func (r *LogCompressionTaskRepositoryImpl) Insert(req model.LogCompressionTask) error {
	_, err := r.WrapDB.MongoLog.Collections.LogCompressionTask.InsertOne(context.TODO(), req)
	if err != nil {
		return err
	}
	return nil
}

func (r *LogCompressionTaskRepositoryImpl) FindById(id string) (res model.LogCompressionTask, err error) {
	filter := bson.M{"task_id": id}

	err = r.WrapDB.MongoLog.Collections.LogCompressionTask.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return
	}

	return
}

func (r *LogCompressionTaskRepositoryImpl) FindAll(f dto.Filtering) (res model.FindAllTaskResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	var filter bson.M
	var records model.LogCompressionTasks

	pgData, err := mongopagination.New(r.WrapDB.MongoLog.Collections.LogCompressionTask).Context(ctx).
		Limit(int64(f.Pagination.Limit)).
		Page(int64(f.Pagination.Page)).
		Sort(f.Sort.SortBy, helper.ConvertSortOrderDirectionForMongo(f.Sort.SortOrder)).
		Filter(filter).
		Decode(&records).Find()

	if err != nil {
		log.Error().Err(err).Msg("err on pgdata")
		return
	}

	res.LogCompressionTasks = records
	res.PaginatedData = pgData

	return
}

func (r *LogCompressionTaskRepositoryImpl) FindAllByJobId(jobId string, f dto.Filtering) (res model.FindAllTaskResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	filter := bson.M{"job_id": jobId}
	var records model.LogCompressionTasks

	pgData, err := mongopagination.New(r.WrapDB.MongoLog.Collections.LogCompressionTask).Context(ctx).
		Limit(int64(f.Pagination.Limit)).
		Page(int64(f.Pagination.Page)).
		Sort(f.Sort.SortBy, helper.ConvertSortOrderDirectionForMongo(f.Sort.SortOrder)).
		Filter(filter).
		Decode(&records).Find()

	if err != nil {
		log.Error().Err(err).Msg("err on pgdata")
		return
	}

	res.LogCompressionTasks = records
	res.PaginatedData = pgData

	return
}

func (r *LogCompressionTaskRepositoryImpl) CountTaskSummary() (res model.TaskSummaries, err error) {
	log.Info().Msg("LogCompressionTaskRepositoryImpl.CountTaskSummary")

	groupStage := bson.D{{
		"$group", bson.D{
			{
				"_id", bson.D{
					{string(lib.StateKey), lib.StateField},
					{string(lib.ActionKey), lib.ActionField},
					{string(lib.StatusKey), lib.StatusField},
				},
			},
			{
				"total", bson.D{{"$sum", 1}},
			},
		},
	}}

	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	cursor, err := r.WrapDB.MongoLog.Collections.LogCompressionTask.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		log.Error().Err(err).Msg("LogCompressionTaskCollection.Aggregate faeild")
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Error().Err(err).Msg("cursor.All failed")
		return
	}

	for _, result := range results {
		data := result["_id"]

		totalStr := fmt.Sprintf("%d", result["total"])
		totalInt, _ := strconv.Atoi(totalStr)

		iter := reflect.ValueOf(data).MapRange()
		taskSummary := model.TaskSummary{
			Total: totalInt,
		}

		for iter.Next() {
			keyStr := fmt.Sprintf("%s", iter.Key().Interface())
			valStr := fmt.Sprintf("%v", iter.Value().Interface())

			if keyStr == model.KeyState {
				taskSummary.State = valStr
			} else if keyStr == model.KeyAction {
				taskSummary.Action = valStr
			} else if keyStr == model.KeyStatus {
				taskSummary.Status = valStr
			}
		}
		res = append(res, taskSummary)
	}

	return
}
