package repository

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	config "github.com/guhkun13/learn-again/multiple-databases/config"
	database "github.com/guhkun13/learn-again/multiple-databases/database"
	"github.com/guhkun13/learn-again/multiple-databases/internal/dto"
	"github.com/guhkun13/learn-again/multiple-databases/internal/model"
	lib "github.com/guhkun13/learn-again/multiple-databases/lib"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogJobExecutionRepositoryImpl struct {
	Env    *config.EnvironmentVariable
	WrapDB *database.WrapDB
}

func NewLogJobExecutionRepositoryImpl(
	env *config.EnvironmentVariable,
	wrapDB *database.WrapDB,
) LogJobExecutionRepository {
	return &LogJobExecutionRepositoryImpl{
		Env:    env,
		WrapDB: wrapDB,
	}
}

func (r *LogJobExecutionRepositoryImpl) Insert(req model.LogJobExecution) error {
	_, err := r.WrapDB.MongoLog.Collections.LogJobExecution.InsertOne(context.TODO(), req)
	if err != nil {
		return err
	}

	return nil
}

func (r *LogJobExecutionRepositoryImpl) CountJobSummaryWithFilter(f dto.LogCountFilter) (int64, error) {
	filter := bson.M{}

	if f.State != "" {
		filter["state"] = f.State
	}
	if f.Action != "" {
		filter["action"] = f.Action
	}

	filter["is_success"] = f.IsSuccess

	count, err := r.WrapDB.MongoLog.Collections.LogJobExecution.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

/** Raw Query on mongo DB */

/*
db.getCollection('log_job_execution').aggregate(
  [
    {
      $group: {
        _id: {
          action: '$action',
          state: '$state',
          success: '$success'
        },
        count: { $sum: 1 }
      }
    }
  ],
  { maxTimeMS: 60000, allowDiskUse: true }
);
*/

// For Admin
func (r *LogJobExecutionRepositoryImpl) CountSummaryJobSteps() (res model.JobSummaries, err error) {
	log.Info().Msg("LogJobExecutionRepositoryImpl.CountJobSummaryV2")

	groupStage := bson.D{{
		"$group", bson.D{
			{
				"_id", bson.D{
					{string(lib.StateKey), lib.StateField},
					{string(lib.ActionKey), lib.ActionField},
					{string(lib.IsSuccessKey), lib.IsSuccessField},
				},
			},
			{
				"total", bson.D{{"$sum", 1}},
			},
		},
	}}

	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	cursor, err := r.WrapDB.MongoLog.Collections.LogJobExecution.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		log.Error().Err(err).Msg("LogJobExecutionCollection.Aggregate faeild")
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
		jobSummary := model.JobSummary{
			Total: totalInt,
		}

		for iter.Next() {
			keyStr := fmt.Sprintf("%s", iter.Key().Interface())
			valStr := fmt.Sprintf("%v", iter.Value().Interface())

			if keyStr == model.KeyState {
				jobSummary.State = valStr
			} else if keyStr == model.KeyAction {
				jobSummary.Action = valStr
			} else if keyStr == model.KeyIsSuccess {
				valBool, err := strconv.ParseBool(valStr)
				if err != nil {
					return nil, err
				}
				jobSummary.IsSuccess = valBool
			}
		}
		res = append(res, jobSummary)
	}

	return
}

/**
For User

[
  {
    $match:

	   {
        user_id:
          "71603138-dd01-4021-8e23-2a0189a07082",
      },
  },
  {
    $group: {
      _id: {
        action: "$action",
        state: "$state",
        is_success: "$is_success",
      },
      count: {
        $sum: 1,
      },
    },
  },
]
**/

// for User (saved for reference)
func (r *LogJobExecutionRepositoryImpl) CountSummaryJobStepsByUserId(userId string) (res model.JobSummaries, err error) {
	log.Info().Msgf("LogJobExecutionRepositoryImpl.CountSummaryJobStepsByUserId (%s)", userId)

	matchStage := bson.D{
		{
			"$match", bson.D{
				{string(lib.UserIdKey), userId},
			},
		},
	}

	groupStage := bson.D{
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{string(lib.StateKey), lib.StateField},
						{string(lib.ActionKey), lib.ActionField},
						{string(lib.IsSuccessKey), lib.IsSuccessField},
					},
				},
				{
					"total", bson.D{{"$sum", 1}},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.Env.Database.Timeout.Read)
	defer cancel()

	cursor, err := r.WrapDB.MongoLog.Collections.LogJobExecution.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		log.Error().Err(err).Msg("LogJobExecutionCollection.Aggregate faeild")
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
		jobSummary := model.JobSummary{
			Total: totalInt,
		}

		for iter.Next() {
			keyStr := fmt.Sprintf("%s", iter.Key().Interface())
			valStr := fmt.Sprintf("%v", iter.Value().Interface())

			if keyStr == model.KeyState {
				jobSummary.State = valStr
			} else if keyStr == model.KeyAction {
				jobSummary.Action = valStr
			} else if keyStr == model.KeyIsSuccess {
				valBool, err := strconv.ParseBool(valStr)
				if err != nil {
					return nil, err
				}
				jobSummary.IsSuccess = valBool
			}
		}
		res = append(res, jobSummary)
	}

	return
}
