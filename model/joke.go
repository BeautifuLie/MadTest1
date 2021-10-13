package model

type Joke struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
	ID    string `json:"id"`
}
