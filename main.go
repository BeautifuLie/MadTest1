package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"program/handlers"
	"program/joker"
	"program/logging"
	"program/storage/mongostorage"
	"program/users"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	logger := logging.InitZapLog()
	godotenv.Load(".env")

	mongoStorage, err := mongostorage.NewMongoStorage(os.Getenv("MONGODB_URI"))
	if err != nil {
		logger.Errorw("Error during connect...", "error", err)
	}

	jokerServer := joker.NewJokerServer(mongoStorage)
	userServer := users.NewUserServer(mongoStorage)

	myRouter := handlers.HandleRequest(handlers.RetHandler(logger, jokerServer, userServer))

	s := http.Server{
		Addr:         ":9090",
		Handler:      myRouter,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.Info(err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	sig := <-signalCh

	logger.Infof("got signal:%", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}

	err = mongoStorage.CloseClientDB()
	if err != nil {
		logger.Info(err)
	}
	logger.Info("Connection to MongoDB closed...")

	logger.Info("Shutdown...")

}
