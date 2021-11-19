package joker

import (
	"errors"
	"math/rand"
	"net/url"
	"program/model"
	"program/storage"
	"strconv"
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

	s.LoadJokesToStruct()

	s.LoadJokesToMap()

	return s
}

//Errors
var ErrNoMatches = errors.New(" No matches")
var ErrLimitOut = errors.New(" Limit out of range")

//Vars

func (s *Server) ID(id string) (model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, err := s.storage.FindID(id)
	if err != nil {
		return model.Joke{}, ErrNoMatches
	}
	return res, nil
}

func (s *Server) Text(text string) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// text = strings.ToLower(strings.TrimSpace(text))
	res, err := s.storage.TextS(text)
	if err != nil {
		return nil, ErrNoMatches
	}
	return res, nil
	// var result []model.Joke

	// text = strings.ToLower(strings.TrimSpace(text))

	// for _, v := range s.jokesStruct {
	// 	title := strings.ToLower(v.Title)
	// 	body := strings.ToLower(v.Body)

	// 	if strings.Contains(title, text) || strings.Contains(body, text) {
	// 		result = append(result, v)
	// 	}
	// }

	// if result != nil {
	// 	return []model.Joke{}, ErrNoMatches
	// }
	// return result, nil
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
	res, err := s.storage.Load()
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

func (s *Server) Add(j model.Joke) (model.Joke, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.storage.Save(j)
	if err != nil {
		return model.Joke{}, errors.New("error writing file")
	}

	return j, nil
}

func (s *Server) LoadJokesToStruct() []model.Joke {

	res, err := s.storage.Load()
	s.jokesStruct = res

	if err != nil {
		return nil
	}
	return s.jokesStruct
}

func (s *Server) LoadJokesToMap() map[string]model.Joke {
	s.jokesMap = map[string]model.Joke{}
	res, err := s.storage.Load()

	if err != nil {
		return nil
	}

	for _, j := range res {
		s.jokesMap[j.ID] = j
	}

	return s.jokesMap
}
