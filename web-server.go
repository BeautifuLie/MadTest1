package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type JokeSource interface {
	LoadJokes() ([]Joke, map[string]Joke)
	SaveJokes([]Joke, map[string]Joke)
}

var a JokeSource

type Joke struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
	ID    string `json:"id"`
}

type Server struct {
	jokesStruct []Joke
	jokesMap    map[string]Joke
}

func main() {

	s := Server{
		[]Joke{},
		map[string]Joke{},
	}

	a = &s
	//jsonUnmarsh(&s)
	myRouter := handleRequest(&s)

	err := http.ListenAndServe(":9090", myRouter)
	if err != nil {
		panic(err)
	}

}

//func jsonUnmarsh(s *Server) {
//
//	j, err := ioutil.ReadFile("reddit_jokes.json")
//	if err != nil {
//		fmt.Println("Error reading file", err)
//	}
//
//	err = json.Unmarshal(j, &s.jokesStruct)
//	if err != nil {
//		fmt.Println("Error unmarshalling JSON", err)
//	}
//
//	sort.SliceStable(s.jokesStruct, func(i, j int) bool {
//		return s.jokesStruct[i].Score > s.jokesStruct[j].Score
//	})
//
//	for _, j := range s.jokesStruct {
//		s.jokesMap[j.ID] = j
//
//	}
//
//}

func (s *Server) LoadJokes() ([]Joke, map[string]Joke) {
	j, err := ioutil.ReadFile("reddit_jokes.json")
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	err = json.Unmarshal(j, &s.jokesStruct)
	if err != nil {
		fmt.Println("Error unmarshalling JSON", err)
	}

	for _, j := range s.jokesStruct {
		s.jokesMap[j.ID] = j
	}
	return s.jokesStruct, s.jokesMap
}
func (s *Server) Load(w http.ResponseWriter, r *http.Request) {
	a.LoadJokes()
	json.NewEncoder(w).Encode("File loaded")
}

func (s *Server) SaveJokes([]Joke, map[string]Joke) {

	structBytes, err := json.MarshalIndent(s.jokesStruct, "", " ")
	if err != nil {
		fmt.Println("Error marshalling JSON", err)
	}
	err = ioutil.WriteFile("jokesStruct.json", structBytes, 0644)

	mapBytes, err := json.MarshalIndent(s.jokesMap, "", " ")
	if err != nil {
		fmt.Println("Error marshalling JSON", err)
	}
	err = ioutil.WriteFile("jokesMap.json", mapBytes, 0644)
}
func (s *Server) Save(w http.ResponseWriter, r *http.Request) {
	a.SaveJokes(s.jokesStruct, s.jokesMap)
	if a.SaveJokes != nil {
		err := json.NewEncoder(w).Encode("File saved")
		if err != nil {
			fmt.Println("Error saving file", err)
		}
	} else {
		json.NewEncoder(w).Encode("Error saving file")
	}
}

func handleRequest(s *Server) *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/jokes", homePage).Methods("GET")
	myRouter.HandleFunc("/jokes/method/save", s.Save)
	myRouter.HandleFunc("/jokes/method/load", s.Load)
	myRouter.HandleFunc("/jokes/funniest", s.getFunniestJokes)
	myRouter.HandleFunc("/jokes/random", s.getRandomJoke)
	myRouter.HandleFunc("/jokes", s.addJoke).Methods("POST")
	myRouter.HandleFunc("/jokes/{id}", s.getJokeByID)
	myRouter.HandleFunc("/jokes/search/{text}", s.getJokeByText)

	return myRouter
}

func (s *Server) getJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	count := 0

	for _, v := range s.jokesMap {

		if strings.Contains(v.ID, id) {
			json.NewEncoder(w).Encode(s.jokesMap[id])
			count++
		}
	}
	if count == 0 {
		http.Error(w, "Error: No jokes found", http.StatusNotFound)
	}

}

func (s *Server) getFunniestJokes(w http.ResponseWriter, r *http.Request) {

	sort.SliceStable(s.jokesStruct, func(i, j int) bool {
		return s.jokesStruct[i].Score > s.jokesStruct[j].Score
	})

	count := 0
	const defaultLimit = 10
	m, _ := url.ParseQuery(r.URL.RawQuery)
	v := ""
	if len(m["limit"]) == 0 {
		count = defaultLimit
	} else {
		v = m["limit"][0] // 0 для того, чтобы брать первый параметр  запроса
		count, _ = strconv.Atoi(v)
	}
	json.NewEncoder(w).Encode(s.jokesStruct[:count])

}

func (s *Server) getRandomJoke(w http.ResponseWriter, r *http.Request) {
	count := 0
	const defaultLimit = 10
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err, "Error parsing query")
	}
	var v string
	if len(m["limit"]) == 0 {
		count = defaultLimit

	} else {
		v = m["limit"][0]
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
	for i := range s.jokesStruct {
		if i < count {
			j := s.jokesStruct[rand.Intn(len(s.jokesStruct))]
			json.NewEncoder(w).Encode(j)

		}

	}

}

func (s *Server) addJoke(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	var j Joke
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.jokesStruct = append(s.jokesStruct, j)
	s.jokesMap[j.ID] = j
	jokeBytes, err := json.Marshal(s.jokesStruct)
	if err != nil {
		fmt.Println("Error marshalling JSON", err)
	}
	err = ioutil.WriteFile("reddit_jokes.json", jokeBytes, 0644)
	if err != nil {
		http.Error(w, "Error creating new Joke", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(&j)
}

func (s *Server) getJokeByText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	text := vars["text"]
	count := 0

	for _, v := range s.jokesStruct {
		fmt.Sprint(count)
		if strings.Contains(v.Title, text) || strings.Contains(v.Body, text) {
			json.NewEncoder(w).Encode(v)

			count++
		}
	}
	if count == 0 {
		http.Error(w, "Error: No matches", http.StatusNotFound)

	}

}
