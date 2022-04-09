package joker

import (
	"fmt"
	"program/model"
	"program/storage"
	"program/tools"
	"strconv"
)

type JokerServer struct {
	jokerServer storage.Joker
}

func NewJokerServer(js storage.Joker) *JokerServer {
	s := &JokerServer{
		jokerServer: js,
	}
	return s
}

func (s *JokerServer) ID(id string) (model.Joke, error) {

	result, err := s.jokerServer.FindID(id)

	if err != nil {
		return model.Joke{}, storage.ErrNoMatches
	}
	return result, nil
}

func (s *JokerServer) Funniest(m string) ([]model.Joke, error) {
	var n int
	a, _ := strconv.Atoi(m)
	if a == 0 {
		n = 10
	} else {
		n = a
	}
	// limit := strconv.Itoa(n)
	result, err := s.jokerServer.Funniest(n)
	if err != nil {
		return nil, fmt.Errorf("funniest jokes error%v", err)
	}
	if len(result) == 0 {

		return nil, storage.ErrNoJokes
	}
	return result, nil
}

func (s *JokerServer) Random(m string) ([]model.Joke, error) {

	var n int
	a, _ := strconv.Atoi(m)
	if a == 0 {
		n = 10
	} else {
		n = a

	}

	res, err := s.jokerServer.Random(n)

	if len(res) == 0 {
		return nil, storage.ErrNoJokes
	}

	if err != nil {
		return nil, fmt.Errorf("random jokes error%v", err)
	}
	if len(res) == 0 {

		return nil, storage.ErrNoJokes
	}

	return res, nil
}

func (s *JokerServer) Text(text string) ([]model.Joke, error) {

	result, err := s.jokerServer.TextSearch(text)
	if err != nil {
		return []model.Joke{}, err
	}

	return result, nil

}
func (s *JokerServer) MonthAndCount(year, count int) (int, int, error) {

	month, count, err := s.jokerServer.MonthAndCount(year, count)
	if err != nil {
		return month, count, err
	}

	return month, count, nil

}
func (s *JokerServer) JokesByMonth(monthNumber int) (int, error) {

	result, err := s.jokerServer.JokesByMonth(monthNumber)
	if err != nil {
		return result, err
	}

	return result, nil

}

func (joker *JokerServer) Add(j model.Joke) (model.Joke, error) {
	_, err := joker.jokerServer.FindID(j.ID)
	if err == storage.ErrNoJokes {
		randTime, _ := tools.RandTimeAndUserID()
		j.Created_at = randTime
		err = joker.jokerServer.AddJoke(j)
		if err != nil {
			return model.Joke{}, err
		}

	} else if err != nil {
		return model.Joke{}, err
	} else {
		return model.Joke{}, fmt.Errorf(" Joke with that ID already exists")
	}

	return j, nil
}

func (s *JokerServer) Update(j model.Joke, id string) (model.Joke, error) {

	err := s.jokerServer.UpdateByID(j.Body, id)
	if err != nil {
		return model.Joke{}, fmt.Errorf("update joke with id %s error:%v", id, err)
	}

	updated, err := s.ID(id)
	if err != nil {
		return model.Joke{}, fmt.Errorf("load joke with id %s error:%v", id, err)
	}

	return updated, nil
}
