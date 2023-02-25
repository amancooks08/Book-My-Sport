package db

import (
	"context"
	"database/sql"
	"fmt"

	logger "github.com/sirupsen/logrus"
)

type User struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Contact  string `db:"contact" json:"contact"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	City     string `db:"city" json:"city"`
	State    string `db:"state" json:"state"`
	Type     string `db:"type" json:"type"`
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
		return err
	}
	return err
}

func (s *pgStore) CheckUser(ctx context.Context, email string, contact string) (bool, error) {
	var flag bool
	err := s.db.QueryRow(CheckUserQuery, &contact, &email).Scan(&flag)
	if err != nil {

		logger.WithField("err", err.Error()).Error("Error checking customer")
		return false, err
	}
	fmt.Printf("flag: %v", flag)
	return flag, err
}

func (s *pgStore) LoginUser(ctx context.Context, email string) (*LoginResponse, error) {
	loginResponse := &LoginResponse{}
	err := s.db.QueryRow(LoginUserQuery, &email).Scan(&loginResponse.Id, &loginResponse.Password, &loginResponse.Role)
	switch {
	case err == sql.ErrNoRows:
		logger.WithField("err", err.Error()).Error("no user with that email id exists.")
		return &LoginResponse{}, err
	case err != nil:
		logger.WithField("err", err.Error()).Error("error logging in customer")
		return &LoginResponse{}, err
	}
	return loginResponse, err
}
