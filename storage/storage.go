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
var ErrNoJokes = errors.New(" No joke in database with such parameters. Create jokes first")

type Storage interface {
	FindID(id string) (model.Joke, error)
	Fun(limit int64) ([]model.Joke, error)
	Random(limit int) ([]model.Joke, error)
	TextSearch(text string) ([]model.Joke, error)
	Save(model.Joke) error
	UpdateByID(text string, id string) (*mongo.UpdateResult, error)
	CloseClientDB() error

	IsExists(model.User) (bool, error)
	CreateUser(user model.User) error
	LoginUser(user model.User) (model.User, error)
	UpdateTokens(signedToken string, signedRefreshToken string, username string) error
}

// type UserStorage interface {
// }
