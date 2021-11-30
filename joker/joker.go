package joker

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"program/model"
	"program/storage"
	"strconv"
	"sync"
)

type Server struct {
	mu      sync.RWMutex
	storage storage.Storage
}

func NewServer(storage storage.Storage) *Server {
	s := &Server{
		storage: storage,
	}

	return s
}

//Errors
var ErrNoMatches = errors.New(" No matches")
var ErrLimitOut = errors.New(" Limit out of range")

func (s *Server) ID(id string) (model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, err := s.storage.FindID(id)
	if err != nil {

		return model.Joke{}, err
	}
	return res, nil
}

func (s *Server) Funniest(m url.Values) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, err := s.storage.Fun()

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
		return nil, ErrLimitOut
	}
	lim := res[:count]
	if lim != nil {
		return lim, nil
	}

	return nil, err
}

func (s *Server) Random(m url.Values) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, err := s.storage.Random()

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
		return nil, ErrLimitOut
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

	return nil, err
}

func (s *Server) Text(text string) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, err := s.storage.TextSearch(text)
	if err != nil {
		return []model.Joke{}, err
	}

	if len(res) == 0 {
		return []model.Joke{}, ErrNoMatches
	}
	return res, nil

}

func (s *Server) Add(j model.Joke) (model.Joke, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.storage.Save(j)
	if err != nil {
		return model.Joke{}, errors.New("error writing file")
	}

	return j, nil
}

func (s *Server) Update(j model.Joke, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := s.storage.UpdateByID(j.Body, id)
	if err != nil {
		return err
	}
	fmt.Print(res)
	return nil
}
