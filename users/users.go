package users

import (
	"program/auth"
	"program/model"
	"program/storage"
	"time"

	"github.com/go-playground/validator"
)

type UserServer struct {
	storage storage.UserStorage
}

func NewUserServer(storage storage.UserStorage) *UserServer {
	s := &UserServer{

		storage: storage,
	}

	return s
}
func (s *UserServer) SignUpUser(u model.User) error {
	var validate = validator.New()
	validationErr := validate.Struct(u)
	if validationErr != nil {
		return storage.ErrPasswordMinLimit
	}
	err := s.storage.IsExists(u)
	if err != nil {
		return err
	}

	password, err := auth.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = password

	u.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	u.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// u.ID = primitive.NewObjectID()

	token, refreshToken, err := auth.GenerateAllTokens(u.Username)
	if err != nil {
		return err
	}
	u.Token = token
	u.Refresh_token = refreshToken
	err = s.storage.CreateUser(u)
	if err != nil {
		return err
	}
	return nil
}
func (s *UserServer) LoginUser(u model.User) (string, error) {
	res, err := s.storage.LoginUser(u)
	if err != nil {
		return "", storage.ErrUserValidate
	}
	err = auth.VerifyPassword(u.Password, res.Password)

	if err != nil {

		return "", err
	}
	token, refreshToken, err := auth.GenerateAllTokens(res.Username)
	if err != nil {
		return "", err
	}

	s.storage.UpdateTokens(token, refreshToken, res.Username)
	return token, nil
}
