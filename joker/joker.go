package joker

import (
	"errors"
	"fmt"
	"log"
	"program/auth"
	"program/model"
	"program/storage"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
		return model.Joke{}, storage.ErrNoMatches
	}
	return result, nil
}

func (s *Server) Funniest(m string) ([]model.Joke, error) {
	var n int64
	a, _ := strconv.Atoi(m)
	if a == 0 {
		n = 10
	} else {
		n = int64(a)
	}

	result, err := s.storage.Fun(n)
	if err != nil {
		return nil, fmt.Errorf("funniest jokes error%v", err)
	}
	if len(result) == 0 {

		return nil, storage.ErrNoJokes
	}
	return result, nil
}

func (s *Server) Random(m string) ([]model.Joke, error) {

	var n int
	a, _ := strconv.Atoi(m)
	if a == 0 {
		n = 10
	} else {
		n = a

	}

	res, err := s.storage.Random(n)

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

func (s *Server) Text(text string) ([]model.Joke, error) {

	result, err := s.storage.TextSearch(text)
	if err != nil {
		return []model.Joke{}, storage.ErrNoMatches
	}

	return result, nil

}

func (s *Server) Add(j model.Joke) (model.Joke, error) {

	err := s.storage.Save(j)
	if err != nil {
		return model.Joke{}, errors.New("error writing file")
	}

	return j, nil
}

func (s *Server) Update(j model.Joke, id string) (model.Joke, error) {

	_, err := s.storage.UpdateByID(j.Body, id)
	if err != nil {
		return model.Joke{}, fmt.Errorf("update joke with id %s error:%v", id, err)
	}

	updated, err := s.ID(id)
	if err != nil {
		return model.Joke{}, fmt.Errorf("load joke with id %s error:%v", id, err)
	}

	return updated, nil
}

var validate = validator.New()

func (s *Server) SignUpUser(u model.User) error {
	validationErr := validate.Struct(u)
	if validationErr != nil {

		return errors.New(" Joker SignUp error")
	}
	ok, err := s.storage.IsExists(u)
	if err != nil {
		return err
	}

	if ok {
		return err
	}

	password := HashPassword(u.Password)
	u.Password = password

	u.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	u.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// u.ID = primitive.NewObjectID()
	// u.User_id = user.ID.Hex()
	token, refreshToken, err := auth.GenerateAllTokens(u.Username)
	if err != nil {
		fmt.Println(err)
	}
	u.Token = token
	u.Refresh_token = refreshToken
	err = s.storage.CreateUser(u)
	if err != nil {
		return err
	}
	return nil
}
func (s *Server) LoginUser(u model.User) (string, error) {
	res, err := s.storage.LoginUser(u)
	if err != nil {
		return "", err
	}
	passwordIsValid, _ := VerifyPassword(u.Password, res.Password)

	if !passwordIsValid {

		return "", errors.New(" Password error")
	}
	token, refreshToken, _ := auth.GenerateAllTokens(res.Username)

	s.storage.UpdateTokens(token, refreshToken, res.Username)
	return token, nil
}

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) //hashsalt
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}
