package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"program/handlers"
	"program/joker"
	"program/logging"
	"program/model"
	"program/storage/mongostorage"
	"program/users"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Init() (*users.UserServer, *joker.JokerServer, *httptest.ResponseRecorder, *zap.SugaredLogger) {
	rr := httptest.NewRecorder()

	logger := logging.InitZapLog()
	mongoStorage, _ := mongostorage.NewMongoStorage("mongodb://localhost:27017")
	s := joker.NewJokerServer(mongoStorage)
	u := users.NewUserServer(mongoStorage)

	return u, s, rr, logger
}
func TestGetJokeByID(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		"/jokes/?id=4xjyho1232", nil)

	user, serv, resp, log := Init()
	h := handlers.RetHandler(log, serv, user)
	handlers.HandleRequest(h)
	h.GetJokeByID(resp, request)
	assert.Equal(t, http.StatusNotFound, resp.Code)

}
func TestGetFunniestJokes(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/funniest?limit=%v", 3), nil)
	user, s, rr, logger := Init()
	h := handlers.RetHandler(logger, s, user)
	handlers.HandleRequest(h)

	h.GetFunniestJokes(rr, request)

	res := rr.Body.Bytes()

	var j []model.Joke

	err := json.Unmarshal(res, &j)
	require.NoError(t, err)

	require.Equal(t, 3, len(j))

	j1 := j[0]
	assert.Equal(t, "On the condition he gets to "+
		"install windows.\n\n\n", j1.Body)
	assert.NotEqual(t, "On the condition he gets to ", j1.Body)
}
func TestGetRandomJoke(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		"/jokes/random", nil)
	user, s, rr, logger := Init()
	h := handlers.RetHandler(logger, s, user)
	handlers.HandleRequest(h)
	h.GetRandomJoke(rr, request)

	///////////////////////////////////////////////////

	request1 := httptest.NewRequest(http.MethodGet,
		"/jokes/random", nil)
	user1, s1, rr1, logger1 := Init()
	h1 := handlers.RetHandler(logger1, s1, user1)
	handlers.HandleRequest(h1)
	h.GetRandomJoke(rr1, request1)

	assert.NotEqual(t, rr, rr1)

}

func TestGetJokeByText(t *testing.T) {
	word := "porcupinetree"

	request := httptest.NewRequest(http.MethodGet,
		"/jokes/search/?text="+word, nil)

	user, s, rr, logger := Init()
	h := handlers.RetHandler(logger, s, user)
	handlers.HandleRequest(h)

	h.GetJokeByText(rr, request)

	assert.Equal(t, 404, rr.Code)

}
