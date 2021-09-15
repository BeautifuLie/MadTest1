package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

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

func main() {
	s := Server{
		jokes:    []Joke{},
		jokesMap: map[string]Joke{},
	}
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
	myRouter.HandleFunc("/jokes/funniest", s.getFunniestJokes)
	myRouter.HandleFunc("/jokes/random", s.getRandomJoke)
	myRouter.HandleFunc("/jokes", s.addJoke).Methods("POST")
	myRouter.HandleFunc("/jokes/{id}", s.getJokeByID)
	myRouter.HandleFunc("/jokes/search/{text}", s.getJokeByText)

	return myRouter
}

//func homePage(w http.ResponseWriter, r *http.Request){
//	tmpl, err := template.ParseFiles("home_page.html")
//	if err != nil {
//		http.Error(w, err.Error(), 400)
//		return
//	}
//
//	err = tmpl.Execute(w, nil)
//	if err != nil {
//		http.Error(w, err.Error(), 400)
//		return
//	}
//}

func (s Server) getJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	json.NewEncoder(w).Encode(s.jokesMap[id])

}

func (s Server) getFunniestJokes(w http.ResponseWriter, r *http.Request) {
	count := 0
	const defaultLimit = 10
	m, _ := url.ParseQuery(r.URL.RawQuery)
	v := ""
	if len(m["limit"]) == 0 {
		count = defaultLimit
	} else {
		v = m["limit"][0]
		count, _ = strconv.Atoi(v)
	}
	json.NewEncoder(w).Encode(s.jokes[:count])

}

func (s Server) getRandomJoke(w http.ResponseWriter, r *http.Request) {
	count := 0
	const defaultLimit = 10
	m, _ := url.ParseQuery(r.URL.RawQuery)
	v := ""
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

func (s Server) addJoke(w http.ResponseWriter, r *http.Request) {

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
	json.NewEncoder(w).Encode(j)
}

func (s Server) getJokeByText(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	text := vars["text"]
	count := 0

	for _, v := range s.jokes {

		if strings.Contains(v.Title, text) || strings.Contains(v.Body, text) {
			json.NewEncoder(w).Encode(v)

			count++
		}
	}
	if count == 0 {
		json.NewEncoder(w).Encode("No matches")
	}

}
