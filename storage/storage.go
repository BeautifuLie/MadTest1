package storage

import (
	"program/model"
)

////go:generate  go run github.com/golang/mock/mockgen -package mocks -destination=./mock_storage.go -source=../storage/storage.go
type Storage interface {
	Load() ([]model.Joke, error)
	Save([]model.Joke) error
}
