package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFunniest(t *testing.T) {
	tt := []struct {
		name    string
		method  string
		limit   int
		wantLen int
	}{
		{
			name:    "limit works",
			method:  http.MethodGet,
			limit:   3,
			wantLen: 3,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			request := httptest.NewRequest(tc.method,
				fmt.Sprintf("/jokes/funniest?limit=%v", tc.limit), nil)
			responseRecorder := httptest.NewRecorder()

			s := Server{
				jokes:    []Joke{},
				jokesMap: map[string]Joke{},
			}

			jsonUnmarsh(&s)
			handleRequest(&s)

			s.getFunniestJokes(responseRecorder, request)

			resp := responseRecorder.Body.Bytes()

			var j []Joke

			json.Unmarshal(resp, &j)

			//assert := assert.New(t)
			//if len(js) != tc.wantLen
			assert.Equal(t, tc.wantLen, len(j),
				fmt.Errorf("Want '%v', got '%v'", tc.wantLen, len(j)))

		})
	}
}

func TestFindById(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		input      string
		want       string
		statusCode int
	}{
		{
			name:       "no ID",
			method:     http.MethodGet,
			input:      "4xjyho1",
			want:       "Error: No jokes found",
			statusCode: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			request := httptest.NewRequest(tc.method,
				fmt.Sprintf("/jokes/%v", tc.input), nil)
			request = mux.SetURLVars(request, map[string]string{"id": tc.input})
			responseRecorder := httptest.NewRecorder()

			s := Server{
				jokes:    []Joke{},
				jokesMap: map[string]Joke{},
			}
			jsonUnmarsh(&s)
			handleRequest(&s)

			s.getJokeByID(responseRecorder, request)

			resp := responseRecorder.Body.Bytes()

			var js []Joke

			json.Unmarshal(resp, &js)

			assert.Equal(t, tc.statusCode, responseRecorder.Code)

		})
	}
}

func TestFindByText(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		input      string
		want       string
		statusCode int
	}{
		{
			name:       "no matches",
			method:     http.MethodGet,
			input:      "porcupinetree",
			want:       "Error: No matches",
			statusCode: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			request := httptest.NewRequest(tc.method,
				fmt.Sprintf("/jokes/search/%v", tc.input), nil)
			request = mux.SetURLVars(request, map[string]string{"text": tc.input})
			responseRecorder := httptest.NewRecorder()

			s := Server{
				jokes:    []Joke{},
				jokesMap: map[string]Joke{},
			}
			jsonUnmarsh(&s)
			handleRequest(&s)

			s.getJokeByText(responseRecorder, request)

			resp := responseRecorder.Body.Bytes()

			var js []Joke

			json.Unmarshal(resp, &js)

			assert.Equal(t, tc.statusCode, responseRecorder.Code)

		})
	}
}

func TestAddJoke(t *testing.T) {

	t.Run("add joke", func(t *testing.T) {
		var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		request := httptest.NewRequest(http.MethodPost,
			fmt.Sprint("/jokes/search/"), bytes.NewBuffer(jsonStr))
		//request = mux.SetURLVars(request, map[string]string{"text": tc.input})
		responseRecorder := httptest.NewRecorder()

		s := Server{
			jokes:    []Joke{},
			jokesMap: map[string]Joke{},
		}
		jsonUnmarsh(&s)
		handleRequest(&s)

		s.addJoke(responseRecorder, request)

		resp := responseRecorder.Body.Bytes()

		var js []Joke

		json.Unmarshal(resp, &js)

		assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	})
}
