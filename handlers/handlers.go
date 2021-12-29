package handlers

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"io"
	"net/http"
	"net/http/httptest"
	"program/joker"
	"program/model"
	"program/storage"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type apiHandler struct {
	Server *joker.Server
	logger *zap.SugaredLogger
}

func RetHandler(logger *zap.SugaredLogger, server *joker.Server) *apiHandler {
	return &apiHandler{
		Server: server,
		logger: logger,
	}
}

func LoggingMiddleware(logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := time.Now()
			logger.Debugw("HTTP request received",
				"url", r.URL)

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, r)

			logger.Debugw("HTTP request processed",
				"url", r.URL,
				"duration", time.Since(s),
				"code", rr.Code,
				// "response", rr.Body.String()
			)

			w.WriteHeader(rr.Code)
			_, err := rr.Body.WriteTo(w)
			if err != nil {
				logger.Error(err)
			}
		})
	})

}

func HandleRequest(h *apiHandler) *mux.Router {
	myRouter := mux.NewRouter()

	myRouter.Use(LoggingMiddleware(h.logger))
	myRouter.HandleFunc("/jokes", h.homePage).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes/funniest", h.GetFunniestJokes).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes/random", h.GetRandomJoke).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes", h.AddJoke).Methods(http.MethodPost)
	myRouter.HandleFunc("/jokes/{id}", h.UpdateJoke).Methods(http.MethodPut)
	myRouter.HandleFunc("/jokes/", h.GetJokeByID).Methods(http.MethodGet)
	myRouter.HandleFunc("/jokes/search/", h.GetJokeByText).Methods(http.MethodGet)

	return myRouter
}
func (h *apiHandler) homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main_page.html")
}

func (h *apiHandler) GetJokeByID(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		json.NewEncoder(w).Encode("id field is empty - type id value")
	}
	res, err := h.Server.ID(id)

	if err != nil {
		h.logger.Infow("GetJokeByID ",
			" id:", id)
		h.respondError(err, w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.logger.Errorw("GetJokeByID encoding error ",
			"error", err,
			"id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *apiHandler) GetFunniestJokes(w http.ResponseWriter, r *http.Request) {

	limit := r.FormValue("limit")

	res, err := h.Server.Funniest(limit)

	if err != nil {

		h.respondError(err, w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(res, "", "  ")
	w.Write(b)

	if err != nil {
		h.logger.Error("GetFunniest encoding error ",
			"error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *apiHandler) GetRandomJoke(w http.ResponseWriter, r *http.Request) {

	limit := r.FormValue("limit")
	res, err := h.Server.Random(limit)
	if err != nil {

		h.respondError(err, w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(res, "", "  ")
	w.Write(b)

	if err != nil {
		h.logger.Errorw("GetRandomJoke encoding error",
			"error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
func (h *apiHandler) GetJokeByText(w http.ResponseWriter, r *http.Request) {

	text := r.URL.Query().Get("text")

	res, err := h.Server.Text(text)
	if err != nil {
		h.logger.Errorw("GetJokeByText error",
			"text", text)
		h.respondError(err, w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(res, "", "  ")
	w.Write(b)

	if err != nil {
		h.logger.Errorw("GetJokeByText encoding error",
			"error", err,
			"text", text)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (h *apiHandler) AddJoke(w http.ResponseWriter, r *http.Request) {
	var j model.Joke
	contentType := r.Header.Get("Content-type")

	if contentType == "application/json" {
		err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&j)
		if err != nil {
			h.logger.Errorw("AddJoke error",
				"error", err,
				"text", r.Body)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
		err := r.ParseForm()
		if err != nil {
			h.logger.Errorw("AddJoke error",
				"error", err,
				"text", r.Body)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		j.Title = r.PostFormValue("title")
		j.Body = r.PostFormValue("body")
		j.Score, _ = strconv.Atoi(r.PostFormValue("score"))
		j.ID = r.PostFormValue("id")

	}

	err := model.Joke.Validate(j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("AddJoke error ",
			"validation error", err)
		_ = json.NewEncoder(w).Encode(serverError{
			Code:        "validation_err",
			Description: err.Error(),
		})

		return
	}
	res, err1 := h.Server.Add(j)
	if err1 != nil {
		h.logger.Errorw("AddJoke error",
			"error", err1)
		h.respondError(err1, w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalIndent(res, "", "  ")
	w.Write(b)

	if err != nil {
		h.logger.Errorw("AddJoke encoding error",
			"error", err,
			"text", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *apiHandler) UpdateJoke(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var j model.Joke

	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&j)
	if err != nil {
		h.logger.Errorw("UpdateJoke error",
			"error", err,
			"text", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(j.Body) == 0 {
		h.logger.Error("UpdateJoke: Body is empty error")
		http.Error(w, "Body is empty", http.StatusBadRequest)
		return

	}

	res, err := h.Server.Update(j, id)
	if err != nil {
		h.logger.Errorw("UpdateJoke  error",
			"error", err)
		h.respondError(err, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalIndent(res, "", "  ")
	w.Write(b)

	if err != nil {
		h.logger.Errorw("UpdateJoke encoding error",
			"error", err,
			"text", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

type serverError struct {
	Code        string
	Description string
}

func (h *apiHandler) respondError(err error, w http.ResponseWriter, r *http.Request) {
	h.logger.Errorw("HTTP respond error",
		"error", err,
		"url", r.URL)
	if errors.Is(err, storage.ErrNoMatches) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if errors.Is(err, storage.ErrLimitOut) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, storage.ErrNoJokes) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
