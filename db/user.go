package db

import (
	"context"
	"database/sql"

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

func (s *pgStore) RegisterUser(ctx context.Context, user *User) error {
	sqlQuery := `INSERT INTO "user" (name, contact, email, password, city, state, type)
    VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`
	err := s.db.QueryRow(sqlQuery, &user.Name, &user.Contact, &user.Email, &user.Password, &user.City, &user.State, &user.Type).Scan(&user.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error registering customer")
		return err
	}
	return err
}

func (s *pgStore) LoginUser(ctx context.Context, email string) (string, error) {
	sqlQuery := `SELECT password FROM "user" WHERE email = $1 `
	var pass string
	err := s.db.QueryRow(sqlQuery, &email).Scan(&pass)
	switch {
	case err == sql.ErrNoRows:
		logger.WithField("err", err.Error()).Error("No user with that Email Id exists.")
		return "", err
	case err != nil:
		logger.WithField("err", err.Error()).Error("Error logging in customer")
		return "", err
	}
	return pass, err
}
