package service

import (
	"context"
	"time"

	"github.com/amancooks08/BookMySport/db"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	logger "github.com/sirupsen/logrus"
)

var secretKey = []byte("secret@987")

type Services interface {
	RegisterUser(ctx context.Context, user *db.User) error
	LoginUser(ctx context.Context, email string, password string) (string, error)
	AddVenue(ctx context.Context, venue *db.Venue) error
	GetAllVenues(ctx context.Context) ([]*db.Venue, error)
	GetVenue(ctx context.Context, name string) (*db.Venue, error)
	UpdateVenue(ctx context.Context, venue *db.Venue, id int) error
	DeleteVenue(ctx context.Context, id int) error
	CheckAvailability(ctx context.Context, id int, date string) ([]*db.Slot, error)
	BookSlot(ctx context.Context, b *db.Booking) error
	GetAllBookings(ctx context.Context, userId int) ([]*db.Booking, error)
	GetBooking(ctx context.Context, bookingid int) (*db.Booking, error)
	CancelBooking(ctx context.Context, id int) error
}

func GenerateToken(loginResponse *db.LoginResponse) (string, error) {
	tokenExpirationTime := time.Now().Add(time.Hour * 24)
	tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": loginResponse.Id,
		"role":    loginResponse.Role,
		"exp":     tokenExpirationTime.Unix(),
	})
	token, err := tokenObject.SignedString(secretKey)
	return token, err
}

type UserOps struct {
	storer db.Storer
}

func NewCustomerOps(storer db.Storer) Services {
	return &UserOps{
		storer: storer,
	}
}

func (cs *UserOps) RegisterUser(ctx context.Context, user *db.User) error {
	user.Password, _ = HashPassword(user.Password)
	err := cs.storer.RegisterUser(ctx, user)
	return err
}

func (cs *UserOps) LoginUser(ctx context.Context, email string, password string) (string, error) {
	loginResponse, err := cs.storer.LoginUser(ctx, email)
	if bcrypt.CompareHashAndPassword([]byte(loginResponse.Password), []byte(password)) != nil {
		return "", err
	}

	token, err := GenerateToken(loginResponse)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error generating fwt token for userId - %s", loginResponse.Id)
	}
	return token, err
}

func (cs *UserOps) AddVenue(ctx context.Context, venue *db.Venue) error {

	err := cs.storer.AddVenue(ctx, venue)
	return err
}

func (cs *UserOps) GetAllVenues(ctx context.Context) ([]*db.Venue, error) {
	venues, err := cs.storer.GetAllVenues(ctx)
	return venues, err
}

func (cs *UserOps) GetVenue(ctx context.Context, name string) (*db.Venue, error) {
	venue, err := cs.storer.GetVenue(ctx, name)
	return venue, err
}

func (cs *UserOps) UpdateVenue(ctx context.Context, venue *db.Venue, id int) error {
	err := cs.storer.UpdateVenue(ctx, venue, id)
	return err
}

func (cs *UserOps) DeleteVenue(ctx context.Context, id int) error {
	err := cs.storer.DeleteVenue(ctx, id)
	return err
}

func (cs *UserOps) CheckAvailability(ctx context.Context, id int, date string) ([]*db.Slot, error) {
	slots, err := cs.storer.CheckAvailability(ctx, id, date)
	return slots, err
}

func (cs *UserOps) BookSlot(ctx context.Context, b *db.Booking) error {
	err := cs.storer.BookSlot(ctx, b)
	return err
}

func (cs *UserOps) GetAllBookings(ctx context.Context, userId int) ([]*db.Booking, error) {
	bookings, err := cs.storer.GetAllBookings(ctx, userId)
	return bookings, err
}

func (cs *UserOps) GetBooking(ctx context.Context, id int) (*db.Booking, error) {
	booking, err := cs.storer.GetBooking(ctx, id)
	return booking, err
}

func (cs *UserOps) CancelBooking(ctx context.Context, id int) error {
	err := cs.storer.CancelBooking(ctx, id)
	return err
}
