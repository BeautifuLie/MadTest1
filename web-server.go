package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"

)

type Joke struct {


	Body  string `json:"body"`
	ID    string `json:"id"`
	Score int   `json:"score"`
	Title string `json:"title"`
	
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

	sort.SliceStable(jokes, func(i, j int) bool { return jokes[i].Score > jokes[j].Score })


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
			json.NewEncoder(w).Encode(v)
			break
		}
	}

}

func getFunniestJoke (w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(jokes[0:10])
}


func getRandomJoke (w http.ResponseWriter, r *http.Request) {

	joke := jokes[rand.Intn(len(jokes))]
	json.NewEncoder(w).Encode(joke)
}






