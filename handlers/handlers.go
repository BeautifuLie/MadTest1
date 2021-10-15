package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"program/joker"
	"program/model"
	"program/storage"
)

type ApiHandler struct {
	Server joker.Server
}

func HandleRequest(h *ApiHandler) *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/jokes", homePage).Methods("GET")
	myRouter.HandleFunc("/jokes/method/save", h.Save)
	myRouter.HandleFunc("/jokes/method/load", h.Load)
	myRouter.HandleFunc("/jokes/funniest", h.GetFunniestJokes)
	myRouter.HandleFunc("/jokes/random", h.GetRandomJoke)
	myRouter.HandleFunc("/jokes", h.AddJoke).Methods("POST")
	myRouter.HandleFunc("/jokes/{id}", h.GetJokeByID)
	myRouter.HandleFunc("/jokes/search/{text}", h.GetJokeByText)

	return myRouter
}

func (h *ApiHandler) Load(w http.ResponseWriter, r *http.Request) {

	_, err := h.Server.JStruct()

	if err != nil {
		switch errors.Cause(err).(type) {
		case *os.PathError:
			w.Write([]byte(err.Error()))

		default:

			w.Write([]byte("other error"))
		}

	} else {
		json.NewEncoder(w).Encode("File loaded")
	}

}

func (h *ApiHandler) GetJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	res, err := h.Server.ID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *ApiHandler) GetJokeByText(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	text := vars["text"]
	res, err := h.Server.Text(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)

}

func (h *ApiHandler) GetFunniestJokes(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}
	res, err := h.Server.Funniest(m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *ApiHandler) GetRandomJoke(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err, "Error parsing query")
	}

	res, err := h.Server.Random(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *ApiHandler) AddJoke(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(serverError{
			Code:        "validation_err",
			Description: err.Error(),
		})
		return
	}
	res, err1 := h.Server.Add(j)
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}

func (h *ApiHandler) Save(w http.ResponseWriter, r *http.Request) {
	str, err := h.Server.JStruct()
	storage.St.Save(str)
	//err := storage.St.Save(h.Server.JStruct())
	if err != nil {
		http.Error(w, "error saving file", 500)
	} else {
		json.NewEncoder(w).Encode("File saved")
	}

}
