package mongostorage

import (
	"context"
	"fmt"
	"program/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

func NewMongoStorage(logger *zap.SugaredLogger, connectURI string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectURI))
	if err != nil {
		logger.Fatalf("error while connecting to mongo: %v", err)
		return nil, fmt.Errorf(" error while connecting to mongo: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		logger.Fatalf("error while pinging mongo: %v", err)
		return nil, fmt.Errorf("pinging mongo: %w", err)
	}

	db := client.Database("mongoData")

	ms := &MongoStorage{
		client:     client,
		collection: db.Collection("Jokes"),
		logger:     logger,
	}

	model := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "body", Value: "text"},
			}},
		{
			Keys: bson.D{
				{Key: "score", Value: -1}},
		},
	}
	_, err = ms.collection.Indexes().CreateMany(context.TODO(), model)
	if err != nil {
		logger.Errorf("error creating indexes: %v", err)
		return nil, err
	}

	return ms, nil
}
func (ms *MongoStorage) FindID(id string) (model.Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var j model.Joke

	err := ms.collection.FindOne(ctx, bson.M{"id": id}).Decode(&j)
	if err != nil {
		ms.logger.Errorw("Storage FindID error:%", "error", err)
		if err == mongo.ErrNoDocuments {
			ms.logger.Debug("Storage FindID error: ", err)
			return model.Joke{}, mongo.ErrNoDocuments
		}
		return model.Joke{}, fmt.Errorf("failed to execute query,error:%w", err)
	}

	return j, nil

}

func (ms *MongoStorage) Fun() ([]model.Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var j []model.Joke

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "score", Value: -1}})

	result, err := ms.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		ms.logger.Error("Storage Fun error:", err)
		return nil, err
	}

	if err = result.All(ctx, &j); err != nil {
		ms.logger.Error("Storage Fun decode error", err)
		return nil, err
	}

	return j, nil
}

func (ms *MongoStorage) Random() ([]model.Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var j []model.Joke

	result, err := ms.collection.Find(ctx, bson.D{})
	if err != nil {
		ms.logger.Error("Storage Random error:", err)
		return nil, err
	}

	if err = result.All(ctx, &j); err != nil {
		ms.logger.Error("Storage Random decode error", err)
		return nil, err
	}

	return j, nil
}

func (ms *MongoStorage) TextSearch(text string) ([]model.Joke, error) {
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

	result, err := ms.collection.Find(ctx, filter)
	if err != nil {
		ms.logger.Error("Storage TextSearch error", err)
		return nil, err
	}

	if err = result.All(ctx, &j); err != nil {
		ms.logger.Error("Storage TextSearch decode error ", err)
		return nil, err
	} else if len(j) == 0 {
		ms.logger.Debug("Storage TextSearch error: ", err)
		return nil, mongo.ErrNoDocuments
	}

	return j, nil

}

func (ms *MongoStorage) Save(j model.Joke) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := ms.collection.InsertOne(ctx, j)
	if err != nil {
		ms.logger.Error("Storage Save error ", err)
		return err
	}
	return nil
}

func (ms *MongoStorage) UpdateByID(text string, id string) (*mongo.UpdateResult, error) {

	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "body", Value: text}}}}

	res, err := ms.collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		ms.logger.Error("Storage UpdateByID error ", err)
		return nil, err
	}

	return res, nil
}

func (ms *MongoStorage) CloseClientDB() {

	if ms.client == nil {
		return
	}

	err := ms.client.Disconnect(context.TODO())
	if err != nil {
		ms.logger.Error("Storage CloseClientDB error ", err)
	}
	ms.logger.Info("Connection to MongoDB closed...")

}
