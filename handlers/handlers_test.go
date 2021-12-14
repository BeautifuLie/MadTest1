package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"program/handlers"
	"program/joker"
	"program/logging"
	"program/model"
	"program/storage/mongostorage"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJokeByID(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		"/jokes/{id}", nil)
	request = mux.SetURLVars(request, map[string]string{"id": "4xjyho1"})
	responseRecorder := httptest.NewRecorder()

	logger := logging.InitZapLog()
	mongoStorage, _ := mongostorage.NewMongoStorage(logger, "mongodb://localhost:27017")
	s := joker.NewServer(logger, mongoStorage)
	h := handlers.RetHandler(logger, s)
	handlers.HandleRequest(h)

	h.GetJokeByID(responseRecorder, request)
	assert.Equal(t, 404, responseRecorder.Code)

}
func TestGetFunniestJokes(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/funniest?limit=%v", 3), nil)
	responseRecorder := httptest.NewRecorder()

	logger := logging.InitZapLog()
	mongoStorage, _ := mongostorage.NewMongoStorage(logger, "mongodb://localhost:27017")

	s := joker.NewServer(logger, mongoStorage)

	h := handlers.RetHandler(logger, s)
	handlers.HandleRequest(h)

	h.GetFunniestJokes(responseRecorder, request)

	resp := responseRecorder.Body.Bytes()

	var j []model.Joke

	err := json.Unmarshal(resp, &j)
	require.NoError(t, err)

	require.Equal(t, 3, len(j))

	j1 := j[0]
	assert.Equal(t, "On the condition he gets to "+
		"install windows.\n\n\n", j1.Body)
}
func TestGetRandomJoke(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		"/jokes/random", nil)
	rr := httptest.NewRecorder()

	logger := logging.InitZapLog()
	mongoStorage, _ := mongostorage.NewMongoStorage(logger, "mongodb://localhost:27017")
	s := joker.NewServer(logger, mongoStorage)
	h := handlers.RetHandler(logger, s)
	handlers.HandleRequest(h)
	h.GetRandomJoke(rr, request)

	///////////////////////////////////////////////////

	request1 := httptest.NewRequest(http.MethodGet,
		"/jokes/random", nil)
	rr1 := httptest.NewRecorder()

	logger1 := logging.InitZapLog()
	mongoStorage1, _ := mongostorage.NewMongoStorage(logger1, "mongodb://localhost:27017")
	s1 := joker.NewServer(logger1, mongoStorage1)
	h1 := handlers.RetHandler(logger1, s1)
	handlers.HandleRequest(h1)
	h.GetRandomJoke(rr1, request1)

	assert.NotEqual(t, rr, rr1)

}

func TestGetJokeByText(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		"/jokes/search/{text}", nil)
	request = mux.SetURLVars(request, map[string]string{"text": "porcupinetree"})
	responseRecorder := httptest.NewRecorder()

	logger := logging.InitZapLog()
	mongoStorage, _ := mongostorage.NewMongoStorage(logger, "mongodb://localhost:27017")
	s := joker.NewServer(logger, mongoStorage)
	h := handlers.RetHandler(logger, s)
	handlers.HandleRequest(h)

	h.GetJokeByText(responseRecorder, request)

	assert.Equal(t, 404, responseRecorder.Code)

}

func TestAddJoke(t *testing.T) {

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast.",
							  "body":"And go away","score":1,"id":"7q6w5e"}`)
	request := httptest.NewRequest(http.MethodPost,
		"/jokes/", bytes.NewBuffer(jsonStr))

	responseRecorder := httptest.NewRecorder()

	logger := logging.InitZapLog()
	mongoStorage, _ := mongostorage.NewMongoStorage(logger, "mongodb://localhost:27017")
	s := joker.NewServer(logger, mongoStorage)
	h := handlers.RetHandler(logger, s)
	handlers.HandleRequest(h)

	h.AddJoke(responseRecorder, request)
	assert.Equal(t, 201, responseRecorder.Code)

}
