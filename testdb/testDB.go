package testdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"program/model"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDB(file string) {
	err := godotenv.Load("/home/denys/go/src/gitlab/maddevices/.env")
	if err != nil {
		fmt.Println("Error during load environments", "error", err)
	}
	connectURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(connectURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("mongo.Connect() ERROR: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()

	colJokes := client.Database("mongoData").Collection("Jokes")

	byteValues, err := ioutil.ReadFile(file)
	if err != nil {

		fmt.Println("ioutil.ReadFile ERROR:", err)
	} else {
		var docs []model.Joke
		err = json.Unmarshal(byteValues, &docs)
		if err != nil {
			fmt.Println("Unmarshal eerror :", err)
			return
		}
		for i := range docs {
			doc := docs[i]
			_, insertErr := colJokes.InsertOne(ctx, doc)
			if insertErr != nil {
				fmt.Println("InsertOne ERROR:", insertErr)
			}
		}
	}
	colUsers := client.Database("mongoData").Collection("Users")
	user := model.User{
		Username:      "Denys",
		Password:      "1234",
		Token:         "abcd",
		Refresh_token: "abcdf",
	}
	_, err = colUsers.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("error insert test user:%w", err)
	}
}
