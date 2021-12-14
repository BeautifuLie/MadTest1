package joker

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"program/model"
	"program/storage"
	"strconv"

	"go.uber.org/zap"
)

type Server struct {
	logger  *zap.SugaredLogger
	storage storage.Storage
}

func NewServer(logger *zap.SugaredLogger, storage storage.Storage) *Server {
	s := &Server{
		logger:  logger,
		storage: storage,
	}

	return s
}

func (s *Server) ID(id string) (model.Joke, error) {

	result, err := s.storage.FindID(id)

	if err != nil {
		return model.Joke{}, storage.ErrNoMatches
	}
	return result, nil
}

func (s *Server) Funniest(m url.Values) ([]model.Joke, error) {

	result, err := s.storage.Fun()

	if len(result) == 0 {
		return nil, storage.ErrNoJokes
	}

	count := 0
	const defaultLimit = 10
	var v string
	var a int

	if len(m["limit"]) > 0 {
		v = m["limit"][0]
		a, _ = strconv.Atoi(v)
	}
	if a > 0 {
		count = a
	} else {
		count = defaultLimit
	}

	if count > len(result) {
		return nil, storage.ErrLimitOut
	}
	lim := result[:count]
	if lim != nil {
		return lim, nil
	}

	return nil, err
}

func (s *Server) Random(m url.Values) ([]model.Joke, error) {

	res, err := s.storage.Random()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, storage.ErrNoJokes
	}

	var result []model.Joke

	count := 0
	const defaultLimit = 10
	var v string
	var a int

	if len(m["limit"]) > 0 {
		v = m["limit"][0]
		a, _ = strconv.Atoi(v)
	}
	if a > 0 {
		count = a
	} else {
		count = defaultLimit
	}

	if count > len(res) {
		return nil, storage.ErrLimitOut
	}

	for i := range res {
		if i < count { //перебирает до указанного "count"
			a := res[rand.Intn(len(res))]
			result = append(result, a)
		}
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("random jokes error%v", err)
}

func (s *Server) Text(text string) ([]model.Joke, error) {

	result, err := s.storage.TextSearch(text)
	if err != nil {
		return []model.Joke{}, storage.ErrNoMatches
	}

	return result, nil

}

func (s *Server) Add(j model.Joke) (model.Joke, error) {

	err := s.storage.Save(j)
	if err != nil {
		return model.Joke{}, errors.New("error writing file")
	}

	return j, nil
}

func (s *Server) Update(j model.Joke, id string) (model.Joke, error) {

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
