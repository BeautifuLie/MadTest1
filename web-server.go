package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/url"
	"program/joker"
	"program/model"
	"program/storage"
	"program/storage/filestorage"
)

func main() {

	storage.St = &filestorage.FileStorage{}

	//file:= filestorage.NewFileStorage("reddit_jokes.json")
	//fmt.Println(file)

	myRouter := handleRequest(&joker.Server{})

	err := http.ListenAndServe(":9090", myRouter)
	if err != nil {
		panic(err)
	}

}

func handleRequest(s *joker.Server) *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/jokes", homePage).Methods("GET")
	myRouter.HandleFunc("/jokes/method/save", Save)
	myRouter.HandleFunc("/jokes/method/load", Load)
	myRouter.HandleFunc("/jokes/funniest", getFunniestJokes)
	myRouter.HandleFunc("/jokes/random", getRandomJoke)
	myRouter.HandleFunc("/jokes", addJoke).Methods("POST")
	myRouter.HandleFunc("/jokes/{id}", getJokeByID)
	myRouter.HandleFunc("/jokes/search/{text}", getJokeByText)

	return myRouter
}

func getJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	res, err := joker.ID(id, &joker.S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func getJokeByText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	text := vars["text"]
	res, err := joker.Text(text, &joker.S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&res)

}

func getFunniestJokes(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}
	res, err := joker.Funniest(m, &joker.S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func getRandomJoke(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err, "Error parsing query")
	}

	res, err := joker.Random(m, &joker.S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func addJoke(w http.ResponseWriter, r *http.Request) {
	type serverError struct {
		Code        string
		Description string
	}

	var j model.Joke
	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//err = model.Joke.Validate(j)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	json.NewEncoder(w).Encode(serverError{
	//		Code:        "validation_err",
	//		Description: err.Error(),
	//	})
	//	return
	//}
	res, err1 := joker.Add(j, &joker.S)
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}

func Load(w http.ResponseWriter, r *http.Request) {

	_, err := storage.St.Load()
	if err != nil {
		http.Error(w, "error", 404)
	} else {
		json.NewEncoder(w).Encode("File loaded")
	}

}

func Save(w http.ResponseWriter, r *http.Request) {

	err := storage.St.Save(joker.S.JokesStruct)
	if err != nil {
		http.Error(w, "error saving file", 500)
	}
	json.NewEncoder(w).Encode("File saved")

}
