package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
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

			var js []Joke

			json.Unmarshal(resp, &js)

			//assert := assert.New(t)
			//if len(js) != tc.wantLen
			assert.Equal(t, tc.wantLen, len(js),
				fmt.Errorf("Want '%v', got '%v'", tc.wantLen, len(js)))

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
			input:      "36q54t3",
			want:       "Error: No jokes found",
			statusCode: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			request := httptest.NewRequest(tc.method,
				fmt.Sprintf("/jokes/%v", tc.input), nil)
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

			assert.NotEqual(t, responseRecorder.Code, tc.statusCode,
				fmt.Errorf("Want status '%d', got '%d'",
					tc.statusCode, responseRecorder.Code))
			//if responseRecorder.Code != tc.statusCode {
			//	fmt.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			//}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				fmt.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}

		})
	}
}
