package model

import (
	"errors"
	"strings"
)

type Joke struct {
	Title string `json:"title" bson:"title"`
	Body  string `json:"body" bson:"body"`
	Score int    `json:"score" bson:"score"`
	ID    string `json:"id" bson:"id"`
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
