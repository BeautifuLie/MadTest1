package storage

import (
	"program/model"

	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/golang/mock/mockgen/model"
)

//go:generate  go run github.com/golang/mock/mockgen -source storage.go -destination mocks/mock_storage.go -package mocks

type Storage interface {
	FindID(id string) (model.Joke, error)
	Fun() ([]model.Joke, error)
	Random() ([]model.Joke, error)
	TextSearch(text string) ([]model.Joke, error)
	Save(model.Joke) error
	UpdateByID(text string, id string) (*mongo.UpdateResult, error)
}
