package filestorage

import (
	"context"
	"fmt"
	"log"
	"program/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoStorage(connectURI string) *MongoStorage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var jokes []model.Joke

	// model := []mongo.IndexModel{
	// 	{
	// 		Keys: bson.D{{"title", "text"}},
	// 	},
	// 	{
	// 		Keys: bson.D{{"body", 1}},
	// 	},
	// }
	// _, err := ms.collection.Indexes().CreateMany(context.TODO(), model)
	// if err != nil {
	// 	panic(err)
	// }

	res, err := ms.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf(" failed to fetch jokes:%w", err)
	}
	err = res.All(ctx, &jokes)
	if err != nil {
		panic(err)
	}

	return jokes, nil
}

func (ms *MongoStorage) Save(j model.Joke) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := ms.collection.InsertOne(ctx, j)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoStorage) FindID(id string) (model.Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var j model.Joke

	err := ms.collection.FindOne(ctx, bson.M{"id": id}).Decode(&j)
	if err != nil {
		return model.Joke{}, err
	}

	return j, nil

}

func (ms *MongoStorage) Fun() ([]model.Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var j []model.Joke

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "score", Value: -1}})
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var j []model.Joke

	filter := bson.D{

		{"$or", bson.A{
			bson.D{{"body", primitive.Regex{Pattern: text, Options: "i"}}},
			bson.D{{"title", primitive.Regex{Pattern: text, Options: "i"}}},
		}},
	}

	cur, err := ms.collection.Find(ctx, filter)

	// filter := bson.D{{"$text", bson.D{{"$search", text}}}} //for indexModel

	cur.All(ctx, &j)
	if err != nil {

		return []model.Joke{}, err
	}
	return j, nil

}
