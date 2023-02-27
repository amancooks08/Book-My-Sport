package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type customerLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Password string `json:"-"`
	Email   string `json:"email"`
	City    string `json:"city"`
	State   string `json:"state"`
	Type    string `json:"type"`
}

type Venue struct {
	Id	  int         `json:"id"`
	Name    string    `json:"name"`
	Contact string    `json:"contact"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	Address string    `json:"address"`
	Email   string    `json:"email"`
	Games   []string  `json:"games"`
}

type Booking struct {
	Id         int       `json:"id"`
	BookedBy   int       `json:"booked_by"`
	BookedAt   int       `json:"booked_at"`
	Time       time.Time `json:"booking_time"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Game       string    `json:"game"`
	AmountPaid float64   `json:"amount"`
}

type Claims struct {
	Email          string
	Password       string
	StandardClaims jwt.RegisteredClaims
	Role           string
}

type LoginResponse struct {
	Token string
	Message string
}
