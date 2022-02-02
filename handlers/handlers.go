package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"io"
	"net/http"
	"net/http/httptest"
	"program/auth"
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
				"response", rr.Body.String(),
			)

			w.WriteHeader(rr.Code)
			_, err := rr.Body.WriteTo(w)
			if err != nil {
				logger.Error(err)
			}
		})
	})

}

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("Authorization")

		if clientToken == "" {
			fmt.Sprintf("Token get error") //http.Error(...)
			return
		}

		claims, err := auth.ValidateToken(clientToken)
		if err != "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Sprintf("Token not valid")
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))

		// http.Redirect(w, r, "/jokes", http.StatusSeeOther)

	})
}

func HandleRequest(h *apiHandler) *mux.Router {
	r := mux.NewRouter()

	r.Use(LoggingMiddleware(h.logger))

	r.HandleFunc("/signup", h.SignUpPage).Methods(http.MethodGet)
	r.HandleFunc("/login", h.LoginPage).Methods(http.MethodGet)
	r.HandleFunc("/signup", h.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)

	// myRouter.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(JwtVerify)
	// apiRouter.Use(MiddlewareValidateAccessToken)
	s.HandleFunc("/api", h.homePage).Methods(http.MethodGet)
	s.HandleFunc("/api/jokes/funniest", h.GetFunniestJokes).Methods(http.MethodGet)
	s.HandleFunc("/api/jokes/random", h.GetRandomJoke).Methods(http.MethodGet)
	s.HandleFunc("/api/jokes", h.AddJoke).Methods(http.MethodPost)
	s.HandleFunc("/api/jokes/{id}", h.UpdateJoke).Methods(http.MethodPut)
	s.HandleFunc("/api/jokes/", h.GetJokeByID).Methods(http.MethodGet)
	s.HandleFunc("/api/jokes/search/", h.GetJokeByText).Methods(http.MethodGet)

	return r
}

// func (h *apiHandler) joinPage(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "join_page.html")
// }

// // func (h *apiHandler) login(w http.ResponseWriter, r *http.Request) {

// // }
func (h *apiHandler) homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main_page.html")
}

func (h *apiHandler) SignUpPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "sign_page.html")
}
func (h *apiHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login_page.html")
}
func (h *apiHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	// username := r.URL.Query().Get("username")
	// password := r.URL.Query().Get("password")
	var u model.User
	contentType := r.Header.Get("Content-type")

	if contentType == "application/json" {
		err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&u)
		if err != nil {
			h.logger.Errorw("AddUser error",
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
		u.Username = r.PostFormValue("username")
		u.Password = r.PostFormValue("password")

	}
	err := model.User.ValidateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("AddUser error ",
			"validation error:", err)
		_ = json.NewEncoder(w).Encode(serverError{
			Code:        "validation_err",
			Description: err.Error(),
		})

		return
	}
	err = h.Server.SignUpUser(u)
	if err != nil {
		h.logger.Errorw("Create user error",
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("User created"))

}

func (h *apiHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u model.User

	u.Username = r.PostFormValue("username")
	u.Password = r.PostFormValue("password")

	resp, err := h.Server.LoginUser(u)
	if err != nil {
		h.logger.Errorw("Login error",
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:9090/auth/api", nil)
	if err != nil {
		fmt.Println("ERRORORO")
		return
	}

	req.Header.Set("Authorization", resp)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR CLIENT DO")
		return
	}
	fmt.Println(res)

	// json.NewEncoder(w).Encode(resp)
	http.Redirect(w, r, "http://localhost:9090/auth/api", 301)

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
			"validation error:", err)
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
