package main

import (
	"log"
	"net/http"
	"program/handlers"
	"program/joker"
	"program/storage/filestorage"
)

func main() {

	fileName := filestorage.NewFileStorage("db/reddit_jokes.json")

	server := joker.NewServer(fileName)

	myRouter := handlers.HandleRequest(handlers.RetHandler(server))

	log.Fatal(http.ListenAndServe(":9090", myRouter))

}
