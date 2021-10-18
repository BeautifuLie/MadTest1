package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"program/handlers"
	"program/model"
	"program/storage"
	"program/storage/filestorage"
	"testing"
)

func TestGetFunniest(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/funniest?limit=%v", 3), nil)
	responseRecorder := httptest.NewRecorder()

	h := handlers.RetHandler()

	_, err := h.Server.JStruct()
	require.NoError(t, err)

	handlers.HandleRequest(h)

	h.GetFunniestJokes(responseRecorder, request)

	resp := responseRecorder.Body.Bytes()

	var j []model.Joke

	err = json.Unmarshal(resp, &j)
	require.NoError(t, err)

	require.Equal(t, 3, len(j))

	j1 := j[0]
	assert.Equal(t, "On the condition he gets to "+
		"install windows.\n\n\n", j1.Body)
}

func TestFindById(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprint("/jokes/{id}"), nil)
	request = mux.SetURLVars(request, map[string]string{"id": "4xjyho1"})
	responseRecorder := httptest.NewRecorder()

	h := handlers.RetHandler()

	_, err := h.Server.JStruct()
	require.NoError(t, err)
	handlers.HandleRequest(h)

	h.GetJokeByID(responseRecorder, request)
	assert.Equal(t, 404, responseRecorder.Code)

}

func TestFindByText(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprint("/jokes/search/{text}"), nil)
	request = mux.SetURLVars(request, map[string]string{"text": "porcupinetree"})
	responseRecorder := httptest.NewRecorder()

	h := handlers.RetHandler()
	storage.St = &filestorage.FileStorage{}
	_, err := h.Server.JStruct()
	require.NoError(t, err)
	handlers.HandleRequest(h)

	h.GetJokeByText(responseRecorder, request)

	assert.Equal(t, 404, responseRecorder.Code)

}

func TestAddJoke(t *testing.T) {

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast.",
							  "body":"And go away","score":1,"id":"7q6w5e"}`)
	request := httptest.NewRequest(http.MethodPost,
		fmt.Sprint("/jokes/"), bytes.NewBuffer(jsonStr))

	responseRecorder := httptest.NewRecorder()

	h := handlers.RetHandler()

	_, err := h.Server.JStruct()
	require.NoError(t, err)
	handlers.HandleRequest(h)

	h.AddJoke(responseRecorder, request)

	assert.Equal(t, 201, responseRecorder.Code)

}

func TestRandom(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/random"), nil)
	rr := httptest.NewRecorder()

	h := handlers.RetHandler()

	_, err := h.Server.JStruct()
	require.NoError(t, err)
	handlers.HandleRequest(h)
	h.GetRandomJoke(rr, request)

	///////////////////////////////////////////////////

	request1 := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/random"), nil)
	rr1 := httptest.NewRecorder()

	h1 := handlers.RetHandler()

	_, err = h.Server.JStruct()
	require.NoError(t, err)
	handlers.HandleRequest(h1)
	h.GetRandomJoke(rr1, request1)

	assert.NotEqual(t, rr, rr1)

}

func TestLoad(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprint("/jokes/method/load"), nil)
	responseRecorder := httptest.NewRecorder()

	h := handlers.RetHandler()
	h.Load(responseRecorder, request)
	resp := responseRecorder.Body.Bytes()

	assert.Equal(t, "\"File loaded\"\n", string(resp))
}
