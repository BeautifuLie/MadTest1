package storage

import (
	"errors"
	"program/model"

	"go.mongodb.org/mongo-driver/mongo"
)

////go:generate  go run github.com/golang/mock/mockgen -source storage.go -destination mocks/mock_storage.go -package mocks
//go:generate  mockgen  -destination=./mock/mock_storage.go -package=mocks . Storage

//Errors
var ErrNoMatches = errors.New(" No matches")
var ErrLimitOut = errors.New(" Limit out of range")
var ErrNoJokes = errors.New(" No jokes in database. Create jokes first")

type Storage interface {
	FindID(id string) (model.Joke, error)
	Fun() ([]model.Joke, error)
	Random() ([]model.Joke, error)
	TextSearch(text string) ([]model.Joke, error)
	Save(model.Joke) error
	UpdateByID(text string, id string) (*mongo.UpdateResult, error)
	CloseClientDB()
}
