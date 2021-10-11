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
	"testing"
)

func TestGetFunniest(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/funniest?limit=%v", 3), nil)
	responseRecorder := httptest.NewRecorder()

	s := Server{
		jokesStruct: []Joke{},
		jokesMap:    map[string]Joke{},
	}
	_, _, err := s.LoadJokes()
	require.NoError(t, err)
	//jsonUnmarsh(&s)
	handleRequest(&s)
	s.getFunniestJokes(responseRecorder, request)

	resp := responseRecorder.Body.Bytes()

	var j []Joke

	json.Unmarshal(resp, &j)
	require.NoError(t, err)

	require.Equal(t, 3, len(j))

	j1 := j[0]
	assert.Equal(t, "On the condition he gets to install windows.\n\n\n", j1.Body)
}

func TestFindById(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprint("/jokes/{id}"), nil)
	request = mux.SetURLVars(request, map[string]string{"id": "4xjyho1"})
	responseRecorder := httptest.NewRecorder()

	s := Server{
		jokesStruct: []Joke{},
		jokesMap:    map[string]Joke{},
	}
	s.LoadJokes()
	//jsonUnmarsh(&s)
	handleRequest(&s)

	s.getJokeByID(responseRecorder, request)

	resp := responseRecorder.Body.Bytes()

	var js map[string]Joke

	json.Unmarshal(resp, &js)

	assert.Equal(t, 404, responseRecorder.Code)

}

func TestFindByText(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprint("/jokes/search/{text}"), nil)
	request = mux.SetURLVars(request, map[string]string{"text": "porcupinetree"})
	responseRecorder := httptest.NewRecorder()

	s := Server{
		jokesStruct: []Joke{},
		jokesMap:    map[string]Joke{},
	}
	s.LoadJokes()
	//jsonUnmarsh(&s)
	handleRequest(&s)

	s.getJokeByText(responseRecorder, request)

	resp := responseRecorder.Body.Bytes()

	var js []Joke

	json.Unmarshal(resp, &js)

	assert.Equal(t, 404, responseRecorder.Code)

}

func TestAddJoke(t *testing.T) {

	var jsonStr = []byte(`[{"title":"Buy cheese and bread for breakfast.",
							  "body":"And go away","score":50,"id":"7q6w5e"}]`)
	request := httptest.NewRequest(http.MethodPost,
		fmt.Sprint("/jokes/"), bytes.NewBuffer(jsonStr))

	responseRecorder := httptest.NewRecorder()

	s := Server{
		jokesStruct: []Joke{},
		jokesMap:    map[string]Joke{},
	}
	s.LoadJokes()
	//jsonUnmarsh(&s)
	handleRequest(&s)

	s.addJoke(responseRecorder, request)

	resp := responseRecorder.Body.Bytes()

	var js []Joke

	json.Unmarshal(resp, &js)

	assert.Equal(t, 201, responseRecorder.Code)

}

func TestRandom(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/random"), nil)
	rr := httptest.NewRecorder()

	s := Server{
		jokesStruct: []Joke{},
		jokesMap:    map[string]Joke{},
	}
	s.LoadJokes()
	//jsonUnmarsh(&s)
	handleRequest(&s)
	s.getRandomJoke(rr, request)

	///////////////////////////////////////////////////

	request1 := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/jokes/random"), nil)
	rr1 := httptest.NewRecorder()

	s1 := Server{
		jokesStruct: []Joke{},
		jokesMap:    map[string]Joke{},
	}
	s.LoadJokes()
	//jsonUnmarsh(&s)
	handleRequest(&s1)
	s.getRandomJoke(rr1, request1)

	assert.NotEqual(t, rr, rr1)

}
