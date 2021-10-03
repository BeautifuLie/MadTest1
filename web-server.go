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
	"os"
	"sort"
	"strconv"
	"strings"
)

type JokeSource interface {
	//LoadJokes() map[string]Joke
	SaveJokes(jokes map[string]Joke)
}

type Joke struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
	ID    string `json:"id"`
}

type Server struct {
	jokes []Joke

	jokesMap map[string]Joke
}

var a JokeSource

func main() {

	s := Server{
		jokes:    []Joke{},
		jokesMap: map[string]Joke{},
	}

	a = &s

	jsonUnmarsh(&s)
	myRouter := handleRequest(&s)

	err := http.ListenAndServe(":9090", myRouter)
	if err != nil {
		panic(err)
	}

}

func jsonUnmarsh(s *Server) {

	j, err := ioutil.ReadFile("reddit_jokes.json")
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	err = json.Unmarshal(j, &s.jokes)
	if err != nil {
		fmt.Println("Error unmarshalling JSON", err)
	}

	sort.SliceStable(s.jokes, func(i, j int) bool {
		return s.jokes[i].Score > s.jokes[j].Score
	})

	for _, j := range s.jokes {
		s.jokesMap[j.ID] = j

	}

}

func handleRequest(s *Server) *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/jokes", homePage).Methods("GET")
	myRouter.HandleFunc("/jokes/method/save", s.Save)
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
	json.NewEncoder(w).Encode(s.jokes[:count])

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
	for i := range s.jokes {
		if i < count {
			j := s.jokes[rand.Intn(len(s.jokes))]
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
	s.jokes = append(s.jokes, j)
	s.jokesMap[j.ID] = j
	jokeBytes, err := json.Marshal(s.jokes)
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

	for _, v := range s.jokes {
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

func (s *Server) SaveJokes(map[string]Joke) {

	file, err := os.Create("test.txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	file.Close()

	jokeBytes, err := json.MarshalIndent(s.jokes, "", " ")
	if err != nil {
		fmt.Println("Error marshalling JSON", err)
	}
	err = ioutil.WriteFile("test.txt", jokeBytes, 0644)

}

func (s *Server) Save(w http.ResponseWriter, r *http.Request) {
	a.SaveJokes(s.jokesMap)
	json.NewEncoder(w).Encode("File saved")

}
