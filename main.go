package main

import (
	"log"
	"net/http"
	"program/handlers"
	"program/joker"
	"program/storage/filestorage"
)

func main() {

	// fileName := filestorage.NewFileStorage("db/reddit_jokes.json")

	mongoStorage := filestorage.NewMongoStorage("mongodb://localhost:27017")

	server := joker.NewServer(mongoStorage)

	myRouter := handlers.HandleRequest(handlers.RetHandler(server))

	log.Fatal(http.ListenAndServe(":9090", myRouter))

}
