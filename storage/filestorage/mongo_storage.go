package filestorage

import (
	"context"
	"fmt"
	"log"
	"program/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoStorage(connectURI string) *MongoStorage {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectURI))
	if err != nil {
		log.Fatalf("Error while connecting to mongo: %v\n", err)
	}
	db := client.Database("mongoData")

	return &MongoStorage{
		client:     client,
		collection: db.Collection("Jokes"),
	}

}

func (ms *MongoStorage) Load() ([]model.Joke, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var jokes []model.Joke
	res, err := ms.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf(" failed to fetch jokes:%w", err)
	}
	res.All(ctx, &jokes)
	if err != nil {
		panic(err)
	}

	return jokes, nil
}

func (ms *MongoStorage) Save(j model.Joke) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := ms.collection.InsertOne(ctx, j)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoStorage) FindID(id string) (model.Joke, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var j model.Joke

	err := ms.collection.FindOne(ctx, bson.M{"id": id}).Decode(&j)
	if err != nil {
		return model.Joke{}, err
	}

	return j, nil

}

func (ms *MongoStorage) Fun() ([]model.Joke, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var j []model.Joke

	opts := options.Find()
	opts.SetSort(bson.D{{"score", -1}})
	sortCursor, err := ms.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err = sortCursor.All(ctx, &j); err != nil {
		return []model.Joke{}, err
	}

	return j, nil
}

func (ms *MongoStorage) TextS(text string) ([]model.Joke, error) {
	panic("not implemented")
}
