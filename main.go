package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"program/awslogic"
	"program/handlers"
	"program/joker"
	"program/logging"
	"program/storage/awsstorage"
	"program/storage/mongostorage"
	"program/storage/sqlstorage"
	"program/users"

	"github.com/joho/godotenv"
)

func main() {
	logger := logging.InitZapLog()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Errorw("Error during load environments", "error", err)
	}
	var jokerServer *joker.JokerServer
	var userServer *users.UserServer

	mongoStorage, err := mongostorage.NewMongoStorage(os.Getenv("MONGODB_URI"))
	if err != nil {
		logger.Errorw("Error during connect...", "error", err)
	} else {
		jokerServer = joker.NewJokerServer(mongoStorage)
		userServer = users.NewUserServer(mongoStorage)
		logger.Info("Connected to MongoDB database")
	}
	// mongoStorage, err := mongostorage.NewMongoStorage(os.Getenv("MONGODB_URI"))
	// if err != nil {
	// 	logger.Errorw("Error during connect...", "error", err)
	// } else {
	// 	jokerServer = joker.NewJokerServer(mongoStorage)
	// 	userServer = users.NewUserServer(mongoStorage)
	// 	logger.Info("Connected to MongoDB database")
	// }sqlStorage, err := sqlstorage.NewSqlStorage(os.Getenv("MYSQL_URI"))
	// if err != nil {
	// 	logger.Errorw("Error during connect SQL database", "error", err)
	// } else {
	// 	jokerServer = joker.NewJokerServer(sqlStorage)
	// 	userServer = users.NewUserServer(sqlStorage)
	// 	logger.Info("Connected to MYSQL database")
	// }
	// mongoStorage, err := mongostorage.NewMongoStorage(os.Getenv("MONGODB_URI"))
	// if err != nil {
	// 	logger.Errorw("Error during connect...", "error", err)
	// } else {
	// 	jokerServer = joker.NewJokerServer(mongoStorage)
	// 	userServer = users.NewUserServer(mongoStorage)
	// 	logger.Info("Connected to MongoDB database")
	// }
	// awsstor, err := awsstorage.NewAwsStorage(
	// 	os.Getenv("AWS_REGION"),
	// 	os.Getenv("AWS_ACCESS_KEY_ID"),
	// 	os.Getenv("AWS_SECRET_ACCESS_KEY"),
	// 	"")

	// awsstor, err := awsstorage.NewAwsStorage(
	// 	os.Getenv("AWS_REGION"),
	// 	os.Getenv("AWS_ACCESS_KEY_ID"),
	// 	os.Getenv("AWS_SECRET_ACCESS_KEY"),
	// 	"")

	// logger := logging.InitZapLog()
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	logger.Errorw("Error during load environments", "error", err)
	// }

	sqlStorage, err := sqlstorage.NewSqlStorage(os.Getenv("MYSQL_URI"))
	if err != nil {
		logger.Errorw("Error during connect SQL database", "error", err)
	} else {
		jokerServer = joker.NewJokerServer(sqlStorage)
		userServer = users.NewUserServer(sqlStorage)
		logger.Info("Connected to MYSQL database")
	}

	awsstor, err := awsstorage.NewAwsStorage(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"")
	if err != nil {
		logger.Errorw("Error during connect to AWS services", "error", err)
	} else {
		logger.Info("Connected to AWS services")
	}

	awsServer := awslogic.NewAwsServer(awsstor)

	myRouter := handlers.HandleRequest(handlers.RetHandler(logger, jokerServer, userServer, awsServer))

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

	err = sqlStorage.CloseClientDB()
	if err != nil {
		logger.Info(err)
	}
	logger.Info("Connection to MYSQL closed...")

	logger.Info("Shutdown...")
}
