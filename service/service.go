package service

import (
	"context"
	"errors"
	"time"

	"github.com/amancooks08/BookMySport/db"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	logger "github.com/sirupsen/logrus"
)

var secretKey = []byte("secret@987")

type Services interface {
	RegisterUser(ctx context.Context, user *db.User) error
	CheckUser(ctx context.Context, email string, contact string) error
	LoginUser(ctx context.Context, email string, password string) (string, error)
	AddVenue(ctx context.Context, venue *db.Venue) error
	CheckVenue(ctx context.Context, name string, contact string, email string) error
	GetAllVenues(ctx context.Context) ([]*db.Venue, error)
	GetVenue(ctx context.Context, id int) (*db.Venue, error)
	UpdateVenue(ctx context.Context, venue *db.Venue, id int) error
	DeleteVenue(ctx context.Context, id int) error
	CheckAvailability(ctx context.Context, id int, date string) ([]*db.Slot, error)
	BookSlot(ctx context.Context, b *db.Booking) (float64, error)
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
	if err != nil {
		return errors.New("error registering user")
	}
	return nil
}

func (cs *UserOps) LoginUser(ctx context.Context, email string, password string) (string, error) {
	loginResponse, err := cs.storer.LoginUser(ctx, email)
	if bcrypt.CompareHashAndPassword([]byte(loginResponse.Password), []byte(password)) != nil {
		return "", errors.New("error: invalid credentials")
	}

	token, err := GenerateToken(loginResponse)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error generating jwt token for given userId")
		return "", errors.New("error: error generating jwt token for given userId")
	}
	return token, nil
}

func (cs *UserOps) AddVenue(ctx context.Context, venue *db.Venue) error {
	err := cs.storer.AddVenue(ctx, venue)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error adding venue")
		return errors.New("error adding venue")
	}
	return nil
}

func (cs *UserOps) GetAllVenues(ctx context.Context) ([]*db.Venue, error) {
	venues, err := cs.storer.GetAllVenues(ctx)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting venues")
		return nil, errors.New("error getting venues")
	}
	return venues, nil
}

func (cs *UserOps) GetVenue(ctx context.Context, id int) (*db.Venue, error) {
	venue, err := cs.storer.GetVenue(ctx, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting venue")
		return nil, errors.New("error getting venue")
	}
	return venue, nil
}

func (cs *UserOps) UpdateVenue(ctx context.Context, venue *db.Venue, id int) error {
	err := cs.storer.UpdateVenue(ctx, venue, id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *UserOps) DeleteVenue(ctx context.Context, id int) error {
	err := cs.storer.DeleteVenue(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *UserOps) CheckAvailability(ctx context.Context, venueId int, date string) ([]*db.Slot, error) {
	slots, err := cs.storer.CheckAvailability(ctx, venueId, date)
	if err != nil {
		return nil, errors.New("error checking availability")
	}
	return slots, nil
}

func (cs *UserOps) BookSlot(ctx context.Context, b *db.Booking) (float64, error) {
	price, err := cs.storer.BookSlot(ctx, b)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error booking slot")
		return 0.0, errors.New("error booking slot")
	}
	return price, nil

}

func (cs *UserOps) GetAllBookings(ctx context.Context, userId int) ([]*db.Booking, error) {
	bookings, err := cs.storer.GetAllBookings(ctx, userId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting bookings")
		return nil, errors.New("error getting bookings")
	}
	if len(bookings) == 0 {
		return nil, errors.New("no bookings found")
	}
	return bookings, nil
}

func (cs *UserOps) GetBooking(ctx context.Context, id int) (*db.Booking, error) {
	booking, err := cs.storer.GetBooking(ctx, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booking")
		return nil, errors.New("error: error getting booking")
	}
	return booking, err
}

func (cs *UserOps) CancelBooking(ctx context.Context, id int) error {
	err := cs.storer.CancelBooking(ctx, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking")
		return errors.New("error: error cancelling booking")
	}
	return nil
}
func (cs *UserOps) CheckUser(ctx context.Context, email string, contact string) error {
	flag, err := cs.storer.CheckUser(ctx, email, contact)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error checking user")
		return errors.New("error checking user")
	}
	if flag {
		return errors.New("error: user already exists")
	}
	return nil
}


func (cs *UserOps) CheckVenue(ctx context.Context, name string, contact string, email string) error {
	flag, err := cs.storer.CheckVenue(ctx, name, contact, email)
	if flag {
		logger.WithField("err", err.Error()).Error("error: venue already exists")
		return errors.New("error: venue already exists")
	}
	return nil 
}