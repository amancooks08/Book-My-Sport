package db

import (
	"context"
	"database/sql"

	logger "github.com/sirupsen/logrus"
)

type User struct {
	Id       int    `json:"id"` 
	Name	 string `json:"name"`
	Contact  string `json:"contact"`
	Email    string `json:"email"`
	Password string `json:"password"`
	City     string `json:"city"`
	State    string `json:"state"`
	Type     string `json:"type"`
}

type LoginResponse struct {
	Id       int    `db:"id" json:"id"`
	Password string `db:"password" json:"password"`
	Role     string `db:"type" json:"type"`
}

func (s *pgStore) RegisterUser(ctx context.Context, user *User) error {
	err := s.db.QueryRow(RegisterUserQuery, &user.Name, &user.Contact, &user.Email, &user.Password, &user.City, &user.State, &user.Type).Scan(&user.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error registering customer")
		return ErrRegisterUser
	}
	return nil
}

func (s *pgStore) CheckUser(ctx context.Context, email string, contact string) (bool, error) {
	var flag bool
	err := s.db.QueryRow(CheckUserQuery, &contact, &email).Scan(&flag)
	if err != nil {

		logger.WithField("err", err.Error()).Error("Error checking customer")
		return false, ErrCheckUser
	}
	return flag, nil
}

func (s *pgStore) LoginUser(ctx context.Context, email string) (*LoginResponse, error) {
	loginResponse := &LoginResponse{}
	err := s.db.QueryRow(LoginUserQuery, &email).Scan(&loginResponse.Id, &loginResponse.Password, &loginResponse.Role)
	switch {
	case err == sql.ErrNoRows:
		logger.WithField("err", err.Error()).Error("no user with that email id exists.")
		return &LoginResponse{}, ErrUserNotExists
	case err != nil:
		logger.WithField("err", err.Error()).Error("error logging in customer")
		return &LoginResponse{}, ErrLogin
	}
	return loginResponse, nil
}
