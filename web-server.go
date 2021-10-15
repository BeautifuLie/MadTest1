package main

import (
	"fmt"
	"net/http"
	"program/handlers"
	"program/storage/filestorage"
)

func main() {

	fileName := filestorage.NewFileStorage("db/reddit_jokes.json")
	fmt.Println(fileName)

	myRouter := handlers.HandleRequest(handlers.RetHandler())

	err := http.ListenAndServe(":9090", myRouter)
	if err != nil {
		panic(err)
	}

}
