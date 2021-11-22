package filestorage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"program/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func init() {
// 	ms.RegisterIndexes()
// }

type MongoStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoStorage(connectURI string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectURI))
	if err != nil {
		// log.Fatalf("Error while connecting to mongo: %v\n", err)
		return nil, fmt.Errorf(" error while connecting to mongo: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("pinging mongo: %w", err)
	}

	db := client.Database("mongoData")

	ms := &MongoStorage{
		client:     client,
		collection: db.Collection("Jokes"),
	}

	model := mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "body", Value: "text"},
		},
	}
	_, err = ms.collection.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}

	return ms, nil
}

func (ms *MongoStorage) RegisterIndexes() {

}

func (ms *MongoStorage) Random() ([]model.Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var jokes []model.Joke

	opts := options.Find()
	opts.SetLimit(50)
	res, err := ms.collection.Find(ctx, bson.D{}, opts)
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

	result := ms.collection.FindOne(ctx, bson.M{"id": id})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return model.Joke{}, mongo.ErrNoDocuments
		}
		return model.Joke{}, fmt.Errorf("failed to execute query,error:%w", result.Err())
	}

	if err := result.Decode(&j); err != nil {
		return j, fmt.Errorf(" failed to decode document,error:%w", err)
	}
	return j, nil

}

func (ms *MongoStorage) Fun() ([]model.Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var j []model.Joke

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "score", Value: -1}})
	opts.SetLimit(100)
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

	// filter := bson.D{		//search without indexes
	// 	{"$or", bson.A{
	// 		bson.D{{"body", primitive.Regex{Pattern: text, Options: "i"}}},
	// 		bson.D{{"title", primitive.Regex{Pattern: text, Options: "i"}}},
	// 	}},
	// }

	filter := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: text}}}} //for indexModel

	cur, err := ms.collection.Find(ctx, filter)
	if err != nil {
		return []model.Joke{}, err
	}

	err = cur.All(ctx, &j)
	if err != nil {
		return []model.Joke{}, err
	}
	return j, nil

}

func (ms *MongoStorage) UpdateByID(text []byte, id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "body", Value: text}}}}
	_, err := ms.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {

		return err
	}

	return nil
}
