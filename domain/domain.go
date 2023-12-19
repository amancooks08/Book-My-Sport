package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Contact  string `json:"contact"`
	Password string `json:"-"`
	Email    string `json:"email"`
	City     string `json:"city"`
	State    string `json:"state"`
	Type     string `json:"type"`
}

type Venue struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	City    string   `json:"city"`
	State   string   `json:"state"`
	Contact string   `json:"contact"`
	Email   string   `json:"email"`
	Opening string   `json:"opening_time"`
	Closing string   `json:"closing_time"`
	Price   float64  `json:"price"`
	Games   []string `json:"games"`
	Rating  float64  `json:"rating"`
	OwnerID int      `json:"-"`
}

type Booking struct {
	ID          int     `json:"id"`
	CustomerID  int     `json:"customer_id"`
	VenueID     int     `json:"venue_id"`
	BookingDate string  `json:"booking_date"`
	BookingTime string  `json:"booking_time"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Game        string  `json:"game"`
	AmountPaid  float64 `json:"amount"`
}

type Slot struct {
	VenueID   int    `json:"venue_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Claims struct {
	Email          string
	Password       string
	StandardClaims jwt.RegisteredClaims
	Role           string
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type Message struct {
	Message string `json:"message"`
}

type BookingResponse struct {
	Message string  `json:"message"`
	Amount  float64 `json:"amount"`
}
