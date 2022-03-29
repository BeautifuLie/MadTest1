package storage

import (
	"errors"
	"program/model"

	"github.com/aws/aws-sdk-go/service/sqs"
)

////go:generate  go run github.com/golang/mock/mockgen -source storage.go -destination mocks/mock_storage.go -package mocks
////go:generate  go run github.com/golang/mock/mockgen -source storage.go -destination mocks/mock_storage.go -package mocks
//go:generate  mockgen  -destination=./mock/mock_storage.go -package=mocks . BaseStorage

//Errors
var ErrNoMatches = errors.New(" No matches")
var ErrLimitOut = errors.New(" Limit out of range")
var ErrNoJokes = errors.New(" No joke in database with such parameters. Create jokes first")
var ErrPasswordMinLimit = errors.New(" The password must contain at least 6 characters")
var ErrUserValidate = errors.New("no user with this login")
var ErrPasswordInvalid = errors.New(" Password is incorrect")
var ErrIncorrectMessage = errors.New("message is not correct")

type BaseStorage interface {
	FindID(id string) (model.Joke, error)
	AddJoke(j model.Joke) error
	UpdateByID(text string, id string) error
}
type ExStorage interface {
	Funniest(limit int) ([]model.Joke, error)
	Random(limit int) ([]model.Joke, error)
	CloseClientDB() error
}
type JokesSearch interface {
	TextSearch(text string) ([]model.Joke, error)
	JokesByMonth(monthNumber int) (int, error)
	MonthAndCount(year, count int) (int, int, error)
}
type UserStorage interface {
	IsExists(model.User) error
	CreateUser(user model.User) error
	LoginUser(user model.User) (model.User, error)
	UpdateTokens(signedToken string, signedRefreshToken string, username string) error
	UsersWithoutJokes() ([]string, error)
}
type AWSfuncs interface {
	UploadTos3(j model.Joke) error
	ReadS3LambdaReport() ([]byte, error)
	GetQueueUrl(queueName string) (string, error)
	SendMsg(id string) (string, error)
	GetMsg() (*sqs.ReceiveMessageOutput, error)
	DeleteMsg(messageHandle string) error
	UploadMessageTos3(filename string) error
}
type Joker interface {
	BaseStorage
	ExStorage
	JokesSearch
}
