package joker

import (
	"errors"
	"math/rand"
	"net/url"
	"program/model"
	"program/storage"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Server struct {
	mu          sync.RWMutex
	storage     storage.Storage
	jokesStruct []model.Joke
	jokesMap    map[string]model.Joke
}

func NewServer(storage storage.Storage) *Server {
	s := &Server{
		storage:     storage,
		jokesStruct: []model.Joke{},
		jokesMap:    map[string]model.Joke{},
	}

	_, err := s.LoadJokesToStruct()
	if err != nil {
		return nil
	}
	_, err = s.LoadJokesToMap()
	if err != nil {
		return nil
	}

	return s
}

//ErrNoMatches
var ErrNoMatches = errors.New(" No matches")
var ErrLimitOut = errors.New(" Limit out of range")

func (s *Server) ID(id string) (model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.jokesMap[id]; ok {
		return s.jokesMap[id], nil
	}

	return model.Joke{}, ErrNoMatches
}

func (s *Server) Text(text string) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []model.Joke

	for _, v := range s.jokesStruct {
		v.Title = strings.ToLower(v.Title)
		v.Body = strings.ToLower(v.Body)
		text = strings.ToLower(text)
		if strings.Contains(v.Title, text) || strings.Contains(v.Body, text) {
			result = append(result, v)
		}
	}

	if result != nil {
		return result, nil
	}
	return nil, ErrNoMatches
}

func (s *Server) Funniest(m url.Values) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	err := errors.New(" Load file first")

	sort.SliceStable(s.jokesStruct, func(i, j int) bool {
		return s.jokesStruct[i].Score > s.jokesStruct[j].Score
	})

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

	if len(s.jokesStruct) == 0 {
		return nil, err
	}
	if count > len(s.jokesStruct) {
		return nil, ErrLimitOut
	}
	res := s.jokesStruct[:count]
	if res != nil {
		return res, nil
	}

	return []model.Joke{}, err
}

func (s *Server) Random(m url.Values) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	err := errors.New(" Load file first")
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
	if count > len(s.jokesStruct) {
		return nil, ErrLimitOut
	}
	for i := range s.jokesStruct {
		if i < count { //перебирает до указанного "count"
			a := s.jokesStruct[rand.Intn(len(s.jokesStruct))]
			result = append(result, a)
		}
	}
	if result != nil {
		return result, nil
	}
	return nil, err
}

func (s *Server) Add(j model.Joke) (model.Joke, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jokesStruct = append(s.jokesStruct, j)
	err := s.storage.Save(s.jokesStruct)
	if err != nil {
		return model.Joke{}, errors.New("error writing file")
	}

	return j, nil
}

func (s *Server) LoadJokesToStruct() ([]model.Joke, error) {

	res, err := s.storage.Load()
	s.jokesStruct = res

	if err != nil {
		return nil, errors.New(" error opening file")
	}
	return s.jokesStruct, nil
}

func (s *Server) LoadJokesToMap() (map[string]model.Joke, error) {
	s.jokesMap = map[string]model.Joke{}
	res, err := s.storage.Load()

	for _, j := range res {
		s.jokesMap[j.ID] = j
	}
	if err != nil {
		return nil, errors.New(" error opening file")
	}

	return s.jokesMap, nil
}
