package handlers

import (
	"encoding/json"
	"errors"

	"io"
	"log"
	"net/http"
	"net/url"
	"program/joker"
	"program/model"

	"github.com/gorilla/mux"
)

type apiHandler struct {
	Server *joker.Server
}

func RetHandler(server *joker.Server) *apiHandler {
	return &apiHandler{
		Server: server,
	}
}

func HandleRequest(h *apiHandler) *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/jokes", h.homePage).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes/funniest", h.GetFunniestJokes).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes/random", h.GetRandomJoke).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes", h.AddJoke).Methods(http.MethodPost)
	myRouter.HandleFunc("/jokes/{id}", h.GetJokeByID).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes/search/{text}", h.GetJokeByText).Methods(http.MethodGet)

	return myRouter
}
func (h *apiHandler) homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main_page.html")
}

func (h *apiHandler) GetJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	res, err := h.Server.ID(id)

	if err != nil {
		respondError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *apiHandler) GetJokeByText(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	text := vars["text"]
	res, err := h.Server.Text(text)
	if err != nil {
		respondError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *apiHandler) GetFunniestJokes(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}

	res, err1 := h.Server.Funniest(m)

	if err1 != nil {
		respondError(err1, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *apiHandler) GetRandomJoke(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err, "Error parsing query")
	}

	res, err1 := h.Server.Random(m)
	if err1 != nil {
		respondError(err1, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *apiHandler) AddJoke(w http.ResponseWriter, r *http.Request) {

	var j model.Joke
	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = model.Joke.Validate(j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(serverError{
			Code:        "validation_err",
			Description: err.Error(),
		})

		if err != nil {
			http.Error(w, "error saving file", http.StatusInternalServerError)
		}
		return
	}
	res, err1 := h.Server.Add(j)
	if err1 != nil {
		respondError(err1, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type serverError struct {
	Code        string
	Description string
}

func respondError(err error, w http.ResponseWriter) {
	if errors.Is(err, joker.ErrNoMatches) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if errors.Is(err, joker.ErrLimitOut) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
