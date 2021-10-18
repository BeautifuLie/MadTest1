package model

import (
	"errors"
	"strings"
)

type Joke struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
	ID    string `json:"id"`
}

func (j Joke) Validate() error {
	if strings.TrimSpace(j.Body) == "" {
		return errors.New("joke Body is empty")
	}
	if strings.TrimSpace(j.ID) == "" {
		return errors.New("ID is empty")
	}
	if strings.TrimSpace(j.Title) == "" {
		return errors.New(" Title is empty")
	}
	if j.Score < 0 {
		return errors.New(" Score is less than 0")
	}
	return nil
}
