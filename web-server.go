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
)

type Joke struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
	ID    string `json:"id"`
}

var jokes []Joke

func main() {

	jsonUnmarsh()
	handleRequest()

}

func jsonUnmarsh() {
	j, err := ioutil.ReadFile("reddit_jokes.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(j, &jokes)
	if err != nil {
		fmt.Println("Error umarshalling JSON", err)
	}

	sort.SliceStable(jokes, func(i, j int) bool {
		return jokes[i].Score > jokes[j].Score
	})

}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/jokes/funniest", getFunniestJoke)
	myRouter.HandleFunc("/jokes/random", getRandomJoke)
	myRouter.HandleFunc("/jokes/{id}", getJokeByID)

	http.ListenAndServe(":9090", myRouter)
}

func getJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	for _, v := range jokes {
		if id == v.ID {
			json.NewEncoder(w).Encode(v.Title)
			json.NewEncoder(w).Encode(v.Body)
			json.NewEncoder(w).Encode(v.ID)
			break
		}
	}

}

func getFunniestJoke(w http.ResponseWriter, r *http.Request) {

	m, _ := url.ParseQuery(r.URL.RawQuery)
	var v = m["limit"][0]
	a, _ := strconv.Atoi(v)

	json.NewEncoder(w).Encode(jokes[:a])

}

func getRandomJoke(w http.ResponseWriter, r *http.Request) {

	for i, _ := range jokes {
		if i < rand.Intn(15)+1 {
			j := jokes[rand.Intn(len(jokes))]
			json.NewEncoder(w).Encode(j.Title)
			json.NewEncoder(w).Encode(j.Body)
			json.NewEncoder(w).Encode("******")
		}

	}

}
