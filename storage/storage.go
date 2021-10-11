package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

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
		return errors.New("joke body is empty")
	}
	return nil
}

func ID(id string, s *Server) (Joke, error) {
	count := 0

	err := errors.New("no jokes with that ID")
	for _, v := range s.JokesMap {

		if strings.Contains(v.ID, id) {
			count++
			return s.JokesMap[id], nil

		}
	}

	return Joke{}, err
}

func Text(text string, s *Server) (Joke, error) {
	count := 0

	err := errors.New(" No matches")

	for _, v := range s.JokesStruct {

		if strings.Contains(v.Title, text) || strings.Contains(v.Body, text) {

			count++
			return v, nil
		}

	}
	return Joke{}, err
}

func Funniest(m url.Values, s *Server) ([]Joke, error) {
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
		count, _ = strconv.Atoi(v)
	}
	if s.JokesStruct[:count] != nil {
		return s.JokesStruct[:count], nil
	}
	return []Joke{}, nil
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
	err1 := errors.New(" Some error")
	return Joke{}, err1
}

func Add(j Joke, s *Server) (Joke, error) {
	s.JokesStruct = append(s.JokesStruct, j)
	s.JokesMap[j.ID] = j
	jokeBytes, err := json.Marshal(s.JokesMap)
	if err != nil {
		errors.New("error marshalling")
	}
	err1 := errors.New("no jokes with that ID")
	err2 := ioutil.WriteFile("reddit_jokes.json", jokeBytes, 0644)
	if err2 != nil {
		return j, err1
	}
	return j, err
}
