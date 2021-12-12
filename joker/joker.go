package joker

import (
	"errors"
	"math/rand"
	"net/url"
	"program/model"
	"program/storage"
	"strconv"

	"go.uber.org/zap"
)

type Server struct {
	logger  *zap.SugaredLogger
	storage storage.Storage
}

func NewServer(logger *zap.SugaredLogger, storage storage.Storage) *Server {
	s := &Server{
		logger:  logger,
		storage: storage,
	}

	return s
}

func (s *Server) ID(id string) (model.Joke, error) {

	result, err := s.storage.FindID(id)

	if err != nil {

		s.logger.With(
			"package", "joker",
			"function", "joker.ID",
			"error", err,
		).Info("get by ID failed")

		return model.Joke{}, storage.ErrNoMatches
	}
	return result, nil
}

func (s *Server) Funniest(m url.Values) ([]model.Joke, error) {

	result, err := s.storage.Fun()

	if len(result) == 0 {

		s.logger.With(
			"package", "joker",
			"function", "joker.Funniest",
			"error", "No jokes in database",
		).Info("get Funniest failed")

		return nil, storage.ErrNoJokes
	}

	count := 0
	const defaultLimit = 10
	var v string
	var a int

	if len(m["limit"]) > 0 {
		v = m["limit"][0]
		a, _ = strconv.Atoi(v)
	}
	if a > 0 {
		count = a
	} else {
		count = defaultLimit
	}

	if count > len(result) {
		s.logger.With(
			"package", "joker",
			"function", "joker.Funniest",
			"error", "Query limit out of range",
		).Info("get Funniest failed")
		return nil, storage.ErrLimitOut
	}
	lim := result[:count]
	if lim != nil {
		return lim, nil
	}

	s.logger.With(
		"package", "joker",
		"function", "joker.Funniest",
		"error", err,
	).Error("get Funniest failed")

	return nil, err
}

func (s *Server) Random(m url.Values) ([]model.Joke, error) {

	res, err := s.storage.Random()
	if len(res) == 0 {
		s.logger.With(
			"package", "joker",
			"function", "joker.Random",
			"error", "No jokes in database",
		).Info("get Random failed")
		return nil, storage.ErrNoJokes
	}

	var result []model.Joke

	count := 0
	const defaultLimit = 10
	var v string
	var a int

	if len(m["limit"]) > 0 {
		v = m["limit"][0]
		a, _ = strconv.Atoi(v)
	}
	if a > 0 {
		count = a
	} else {
		count = defaultLimit
	}

	if count > len(res) {
		s.logger.With(
			"package", "joker",
			"function", "joker.Random",
			"error", "Query limit out of range",
		).Info("get Random failed")
		return nil, storage.ErrLimitOut
	}

	for i := range res {
		if i < count { //перебирает до указанного "count"
			a := res[rand.Intn(len(res))]
			result = append(result, a)
		}
	}

	if result != nil {
		return result, nil
	}

	s.logger.With(
		"package", "joker",
		"function", "joker.Random",
		"error", err,
	).Error("get Random failed")

	return nil, err
}

func (s *Server) Text(text string) ([]model.Joke, error) {

	result, err := s.storage.TextSearch(text)
	if err != nil {
		s.logger.With(
			"package", "joker",
			"function", "joker.Text",
			"error", err,
		).Info("get by Text failed")
		return []model.Joke{}, storage.ErrNoMatches
	}

	return result, nil

}

func (s *Server) Add(j model.Joke) (model.Joke, error) {

	err := s.storage.Save(j)
	if err != nil {
		s.logger.With(
			"package", "joker",
			"function", "joker.Add",
			"error", err,
		).Info("Add joke failed")
		return model.Joke{}, errors.New("error writing file")
	}

	return j, nil
}

func (s *Server) Update(j model.Joke, id string) (model.Joke, error) {

	_, err := s.storage.UpdateByID(j.Body, id)
	if err != nil {
		s.logger.With(
			"package", "joker",
			"function", "joker.Update",
			"error", err,
		).Info("Update joke failed")
		return model.Joke{}, err
	}
	updated, err := s.ID(id)
	if err != nil {
		s.logger.With(
			"package", "joker",
			"function", "joker.Update",
			"error", err,
		).Info("Updated joke failed")
		return model.Joke{}, storage.ErrNoJokes
	}

	return updated, nil
}
