package main

import (
	"fmt"
	"net/http"
	"program/handlers"
	"program/storage/filestorage"
)

func main() {

	myRouter := handlers.HandleRequest(handlers.RetHandler())
	fileName := filestorage.NewFileStorage("jokes.json")
	fmt.Println(fileName)

	err := http.ListenAndServe(":9090", myRouter)
	if err != nil {
		panic(err)
	}

}
