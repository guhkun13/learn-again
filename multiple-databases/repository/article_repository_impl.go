package repository

import (
	"context"
	"time"

	"github.com/guhkun13/learn-again/multiple-databases/dto"
	"github.com/guhkun13/learn-again/multiple-databases/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(req dto.CreateArticle) (*mongo.InsertOneResult, error) {
	article := model.Article{}
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	return articleCollection.InsertOne(context.TODO(), article)
}

func getArticle(id primitive.ObjectID) (*Article, error) {
	var article Article
	err := articleCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&article)
	return &article, err
}

func updateArticle(id primitive.ObjectID, update Article) (*mongo.UpdateResult, error) {
	update.UpdatedAt = time.Now()
	return articleCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": update})
}

func deleteArticle(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return articleCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
}

func (r *ArticleRepositoryImpl) Create(req dto.CreateUser) error {
	panic("implement this")
}

func (r *ArticleRepositoryImpl) List() (model.Users, error) {
	panic("implement this")

}

func (r *ArticleRepositoryImpl) Read(id int) (model.User, error) {
	panic("implement this")

}
