package db

import (
	"context"
	"database/sql"

	logger "github.com/sirupsen/logrus"
)

type User struct {
	ID       int    
	Name     string 
	Contact  string 
	Email    string 
	Password string 
	City     string 
	State    string 
	Type     string 
}

type LoginResponse struct {
	Id       int    
	Password string 
	Role     string 
}

func (s *pgStore) RegisterUser(ctx context.Context, user *User) error {
	err := s.db.QueryRow(RegisterUserQuery, &user.Name, &user.Contact, &user.Email, &user.Password, &user.City, &user.State, &user.Type).Scan(&user.ID)
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
	var loginResponse = &LoginResponse{}
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
