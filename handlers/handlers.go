package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"program/joker"
	"program/model"
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
	myRouter.HandleFunc("/jokes", h.homePage).Methods("GET")

	myRouter.HandleFunc("/jokes/funniest", h.GetFunniestJokes)
	myRouter.HandleFunc("/jokes/random", h.GetRandomJoke)
	myRouter.HandleFunc("/jokes", h.AddJoke).Methods("POST")
	myRouter.HandleFunc("/jokes/{id}", h.GetJokeByID)
	myRouter.HandleFunc("/jokes/search/{text}", h.GetJokeByText)

	return myRouter
}
func (h *apiHandler) homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("main_page.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func (h *apiHandler) GetJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	res, err := h.Server.ID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (h *apiHandler) GetJokeByText(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	text := vars["text"]
	res, err := h.Server.Text(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

}

func (h *apiHandler) GetFunniestJokes(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}
	res, err1 := h.Server.Funniest(m)

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

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
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

func (h *apiHandler) AddJoke(w http.ResponseWriter, r *http.Request) {
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

	err = model.Joke.Validate(j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(serverError{
			Code:        "validation_err",
			Description: err.Error(),
		})
		if err != nil {
			http.Error(w, "error saving file", 500)
		}
		return
	}
	res, err1 := h.Server.Add(j)
	if err1 != nil {
		http.Error(w, "Error adding joke", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
