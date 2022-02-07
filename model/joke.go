package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Joke struct {
	Title string `json:"title" bson:"title"`
	Body  string `json:"body" bson:"body"`
	Score int    `json:"score" bson:"score"`
	ID    string `json:"id" bson:"id"`
}
type User struct {
	// ID       primitive.ObjectID `bson:_id`
	Username      string    `json:"username" bson:"username"`
	Password      string    `json:"password" bson:"password" validate:"required,min=6"`
	Token         string    `json:"token" bson:"token"`
	Refresh_token string    `json:"refresh_token" bson:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

func (j Joke) Validate() error {
	if strings.TrimSpace(j.Title) == "" {
		return errors.New(" Title is empty")
	}
	if strings.TrimSpace(j.Body) == "" {
		return errors.New("joke Body is empty")
	}
	if j.Score < 0 {
		return errors.New(" Score is less than 0")
	}
	if strings.TrimSpace(j.ID) == "" {
		return errors.New(" ID is empty")
	}

	return nil
}
func (u User) ValidateUser() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New(" Username is empty")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New(" Password is empty")
	}

	return nil
}
func (user *User) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("error hash")
	}

	return string(bytes)

}

// CheckPassword checks user password
func (user *User) CheckPassword(userPassword string, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	if err != nil {
		check = false
		return check, err
	}

	return check, nil
}
