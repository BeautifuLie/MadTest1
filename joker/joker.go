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
)

type Server struct {
	Storage     storage.Storage
	JokesStruct []model.Joke
	JokesMap    map[string]model.Joke
}

var S = Server{
	Storage:     storage.St,
	JokesStruct: []model.Joke{},
	JokesMap:    map[string]model.Joke{},
}

//ErrNoMatches
var ErrNoMatches = errors.New(" No matches")

func ID(id string, s *Server) (model.Joke, error) {
	for _, j := range S.JokesStruct {
		S.JokesMap[j.ID] = j
	}
	err := errors.New("no jokes with that ID")
	for _, v := range s.JokesMap {

		if strings.Contains(v.ID, id) {
			return s.JokesMap[id], nil
		}
	}
	return model.Joke{}, err
}

func Text(text string, s *Server) ([]model.Joke, error) {

	var result []model.Joke

	for _, v := range s.JokesStruct {
		if strings.Contains(v.Title, text) || strings.Contains(v.Body, text) {
			result = append(result, v)
		}
	}
	if result != nil {
		return result, nil
	}
	return nil, ErrNoMatches
}

func Funniest(m url.Values, s *Server) ([]model.Joke, error) {

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
	return []model.Joke{}, err
}

func Random(m url.Values, s *Server) ([]model.Joke, error) {
	var result []model.Joke
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

			a := S.JokesStruct[rand.Intn(len(s.JokesStruct))]
			result = append(result, a)

		}

	}
	if result != nil {
		return result, nil
	}
	return nil, err
}

func Add(j model.Joke, s *Server) (model.Joke, error) {
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
	err := storage.St.Save(s.JokesStruct)
	if err != nil {
		return model.Joke{}, errors.New("error writing file")
	}

	return j, err
}

//type Joke struct{
//	*model.Joke
//}
