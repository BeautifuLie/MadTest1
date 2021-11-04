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

	s.LoadJokesToStruct()

	s.LoadJokesToMap()

	return s
}

//Errors
var ErrNoMatches = errors.New(" No matches")
var ErrLimitOut = errors.New(" Limit out of range")
var ErrOpenFile = errors.New(" The system cannot find the file")
var ErrNoFile = errors.New(" No file to write")

func (s *Server) ID(id string) (model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.jokesMap) == 0 {
		return model.Joke{}, ErrOpenFile
	}

	if _, ok := s.jokesMap[id]; ok {
		return s.jokesMap[id], nil
	}

	return model.Joke{}, ErrNoMatches
}

func (s *Server) Text(text string) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.jokesStruct) == 0 {
		return nil, ErrOpenFile
	}

	var result []model.Joke

	text = strings.ToLower(strings.TrimSpace(text))

	for _, v := range s.jokesStruct {
		title := strings.ToLower(v.Title)
		body := strings.ToLower(v.Body)

		if strings.Contains(title, text) || strings.Contains(body, text) {
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

	if len(s.jokesStruct) == 0 {
		return nil, ErrOpenFile
	}

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

	if count > len(s.jokesStruct) {
		return nil, ErrLimitOut
	}
	res := s.jokesStruct[:count]
	if res != nil {
		return res, nil
	}

	return nil, ErrOpenFile
}

func (s *Server) Random(m url.Values) ([]model.Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.jokesStruct) == 0 {
		return nil, ErrOpenFile
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

	return nil, ErrOpenFile
}

func (s *Server) Add(j model.Joke) (model.Joke, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.jokesStruct) == 0 {
		return model.Joke{}, ErrNoFile
	}
	s.jokesStruct = append(s.jokesStruct, j)
	err := s.storage.Save(s.jokesStruct)
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
