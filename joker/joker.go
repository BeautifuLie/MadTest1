package joker

import (
	"errors"
	"fmt"
	"program/model"
	"program/storage"
	"strconv"
)

type JokerServer struct {
	storage storage.Storage
}

func NewJokerServer(storage storage.Storage) *JokerServer {
	s := &JokerServer{

		storage: storage,
	}

	return s
}

func (s *JokerServer) ID(id string) (model.Joke, error) {

	result, err := s.storage.FindID(id)

	if err != nil {
		return model.Joke{}, storage.ErrNoMatches
	}
	return result, nil
}

func (s *JokerServer) Funniest(m string) ([]model.Joke, error) {
	var n int64
	a, _ := strconv.Atoi(m)
	if a == 0 {
		n = 10
	} else {
		n = int64(a)
	}

	result, err := s.storage.Fun(n)
	if err != nil {
		return nil, fmt.Errorf("funniest jokes error%v", err)
	}
	if len(result) == 0 {

		return nil, storage.ErrNoJokes
	}
	return result, nil
}

func (s *JokerServer) Random(m string) ([]model.Joke, error) {

	var n int
	a, _ := strconv.Atoi(m)
	if a == 0 {
		n = 10
	} else {
		n = a

	}

	res, err := s.storage.Random(n)

	if len(res) == 0 {
		return nil, storage.ErrNoJokes
	}

	if err != nil {
		return nil, fmt.Errorf("random jokes error%v", err)
	}
	if len(res) == 0 {

		return nil, storage.ErrNoJokes
	}

	return res, nil
}

func (s *JokerServer) Text(text string) ([]model.Joke, error) {

	result, err := s.storage.TextSearch(text)
	if err != nil {
		return []model.Joke{}, storage.ErrNoMatches
	}

	return result, nil

}

func (s *JokerServer) Add(j model.Joke) (model.Joke, error) {

	err := s.storage.Save(j)
	if err != nil {
		return model.Joke{}, errors.New("error writing file")
	}

	return j, nil
}

func (s *JokerServer) Update(j model.Joke, id string) (model.Joke, error) {

	_, err := s.storage.UpdateByID(j.Body, id)
	if err != nil {
		return model.Joke{}, fmt.Errorf("update joke with id %s error:%v", id, err)
	}

	updated, err := s.ID(id)
	if err != nil {
		return model.Joke{}, fmt.Errorf("load joke with id %s error:%v", id, err)
	}

	return updated, nil
}
