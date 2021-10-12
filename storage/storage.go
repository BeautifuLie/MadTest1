package storage

import (
	"errors"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

//go:generate  go run github.com/golang/mock/mockgen -package mocks -destination=./mock_storage.go -source=../storage/storage.go
type Storage interface {
	Load() ([]Joke, error)
	Save([]Joke) error
}

var St Storage

type Joke struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
	ID    string `json:"id"`
}

type Server struct {
	Storage     Storage
	JokesStruct []Joke
	JokesMap    map[string]Joke
}

var S = Server{
	Storage:     St,
	JokesStruct: []Joke{},
	JokesMap:    map[string]Joke{},
}

func (j Joke) Validate() error {
	if strings.TrimSpace(j.Body) == "" {
		return errors.New("joke Body is empty")
	}
	if strings.TrimSpace(j.ID) == "" {
		return errors.New("ID is empty")
	}
	if strings.TrimSpace(j.Title) == "" {
		return errors.New(" Title is empty")
	}
	if j.Score < 0 {
		return errors.New(" Score is less than 0")
	}
	return nil
}

func ID(id string, s *Server) (Joke, error) {

	err := errors.New("no jokes with that ID")
	for _, v := range s.JokesMap {

		if strings.Contains(v.ID, id) {
			return s.JokesMap[id], nil
		}
	}
	return Joke{}, err
}

func Text(text string, s *Server) (Joke, error) {

	err := errors.New(" No matches")

	for _, v := range s.JokesStruct {
		if strings.Contains(v.Title, text) || strings.Contains(v.Body, text) {
			return v, nil
		}
	}

	return Joke{}, err
}

func Funniest(m url.Values, s *Server) ([]Joke, error) {
	err := errors.New(" Bad request")

	sort.SliceStable(s.JokesStruct, func(i, j int) bool {
		return s.JokesStruct[i].Score > s.JokesStruct[j].Score
	})

	count := 0
	const defaultLimit = 10
	v := ""

	if len(m["limit"]) == 0 {
		count = defaultLimit
	} else {
		v = m["limit"][0] // 0 для того, чтобы брать первый параметр  запроса
		count, err = strconv.Atoi(v)
		if err != nil {
			errors.New(" Error converting string to int")
		}
		if count > len(s.JokesStruct) {
			return nil, errors.New(" Limit out of range")
		}

	}
	if s.JokesStruct[:count] != nil {
		return s.JokesStruct[:count], nil
	}
	return []Joke{}, err
}

func Random(m url.Values, s *Server) (Joke, error) {

	count := 0
	const defaultLimit = 10
	var v string
	if len(m["limit"]) == 0 {
		count = defaultLimit
	} else {
		v = m["limit"][0] // 0 для того, чтобы брать первый параметр  запроса
	}
	a, err := strconv.Atoi(v)
	if err != nil {
		println(err.Error())
	}

	if a > 0 {
		count = a
	} else {
		count = defaultLimit
	}
	for i := range s.JokesStruct {

		if i < count {

			return S.JokesStruct[rand.Intn(len(s.JokesStruct))], nil

		}

	}
	err = errors.New(" Some error")
	return Joke{}, err
}

func Add(j Joke, s *Server) (Joke, error) {
	s.JokesStruct = append(s.JokesStruct, j)

	//jokeBytes, err := json.Marshal(s.JokesStruct)
	//if err != nil {
	//	errors.New("error marshalling")
	//}
	//
	//err = ioutil.WriteFile("reddit_jokes.json", jokeBytes, 0644)
	//if err != nil {
	//	return Joke{}, errors.New("error writing file")
	//}
	//return j, err
	err := St.Save(s.JokesStruct)
	if err != nil {
		return Joke{}, errors.New("error writing file")
	}

	return j, err
}
