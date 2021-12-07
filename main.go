package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"program/handlers"
	"program/joker"
	"program/storage/mongostorage"
	"syscall"
	"time"
)

func main() {

	mongoStorage, err := mongostorage.NewMongoStorage("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	server := joker.NewServer(mongoStorage)

	myRouter := handlers.HandleRequest(handlers.RetHandler(server))

	s := http.Server{
		Addr:         ":9090",
		Handler:      myRouter,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		s.ListenAndServe()
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	sig := <-signalCh
	log.Println("got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx)

	mongoStorage.CloseClientDB()

	log.Fatal("shutdown...")

}
