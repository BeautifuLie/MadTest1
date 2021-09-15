package main

import (
	"encoding/json"
	"fmt"
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
			request := httptest.NewRequest(tc.method, fmt.Sprintf("/jokes/funniest?limit=%v", tc.limit), nil)
			responseRecorder := httptest.NewRecorder()

			s := Server{
				jokes:    []Joke{},
				jokesMap: map[string]Joke{},
			}
			jsonUnmarsh(&s)
			handleRequest(&s)

			s.getFunniestJokes(responseRecorder, request)

			resp := responseRecorder.Body.Bytes()

			js := []Joke{}

			json.Unmarshal(resp, &js)

			if len(js) != tc.wantLen {
				t.Errorf("Want '%v', got '%v'", tc.wantLen, len(js))
			}
		})
	}
}
