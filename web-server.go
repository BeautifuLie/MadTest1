package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/url"
	"program/storage"
)

func main() {

	storage.St = &storage.S
	myRouter := handleRequest(&storage.Server{})

	err := http.ListenAndServe(":9090", myRouter)
	if err != nil {
		panic(err)
	}

}

func handleRequest(s *storage.Server) *mux.Router {
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
	res, err := storage.ID(id, &storage.S)
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
	res, err := storage.Text(text, &storage.S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func getFunniestJokes(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}
	res, err1 := storage.Funniest(m, &storage.S)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func getRandomJoke(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err, "Error parsing query")
	}

	res, err1 := storage.Random(m, &storage.S)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func addJoke(w http.ResponseWriter, r *http.Request) {
	type serverError struct {
		Code        string
		Description string
	}

	var j storage.Joke
	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.Joke.Validate(j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serverError{
			Code:        "validation_err",
			Description: err.Error(),
		})
		return
	}
	res, err1 := storage.Add(j, &storage.S)
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
		http.Error(w, "Error loading file", 402)
	} else {
		json.NewEncoder(w).Encode("File loaded")
	}

}

func Save(w http.ResponseWriter, r *http.Request) {

	err := storage.St.Save(storage.S.JokesStruct)
	if err != nil {
		http.Error(w, "oops", 400)
	}
	json.NewEncoder(w).Encode("File saved")

}
