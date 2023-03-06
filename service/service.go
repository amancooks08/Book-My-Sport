package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/amancooks08/BookMySport/db"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	logger "github.com/sirupsen/logrus"

	domain "github.com/amancooks08/BookMySport/domain"
)

var secretKey = []byte("secret@987")

type Services interface {
	RegisterUser(ctx context.Context, user domain.User) error
	CheckUser(ctx context.Context, email string, contact string) error
	LoginUser(ctx context.Context, email string, password string) (string, error)
	AddVenue(ctx context.Context, venue domain.Venue) error
	CheckVenue(ctx context.Context, name string, contact string, email string) error
	GetAllVenues(ctx context.Context) ([]domain.Venue, error)
	GetVenue(ctx context.Context, id int) (domain.Venue, error)
	UpdateVenue(ctx context.Context, venue domain.Venue, userID int, id int) error
	DeleteVenue(ctx context.Context, userID, id int) error
	CheckAvailability(ctx context.Context, id int, date string) ([]domain.Slot, error)
	BookSlot(ctx context.Context, b domain.Booking) (float64, error)
	GetAllBookings(ctx context.Context, userId int) ([]domain.Booking, error)
	GetBooking(ctx context.Context, bookingid int) (domain.Booking, error)
	CancelBooking(ctx context.Context, id int) error
}

func GenerateToken(loginResponse db.LoginResponse) (string, error) {
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

func (cs *UserOps) RegisterUser(ctx context.Context, user domain.User) error {
	dbUser := db.User{
		Name:     user.Name,
		Contact:  user.Contact,
		Email:    user.Email,
		Password: user.Password,
		City:     user.City,
		State:    user.State,
		Type:     user.Type,
	}

	dbUser.Password, _ = HashPassword(user.Password)
	err := cs.storer.RegisterUser(ctx, dbUser)
	if err != nil {
		return errors.New("error registering user")
	}
	return nil
}

func (cs *UserOps) LoginUser(ctx context.Context, email string, password string) (string, error) {
	loginResponse, err := cs.storer.LoginUser(ctx, email)
	if bcrypt.CompareHashAndPassword([]byte(loginResponse.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateToken(loginResponse)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error generating jwt token for given userId")
		return "", errors.New("error generating jwt token for given userId")
	}
	return token, nil
}

func (cs *UserOps) AddVenue(ctx context.Context, venue domain.Venue) error {
	dbVenue := db.Venue{
		Name:    venue.Name,
		Address: venue.Address,
		City:    venue.City,
		State:   venue.State,
		Contact: venue.Contact,
		Email:   venue.Email,
		Opening: venue.Opening,
		Closing: venue.Closing,
		Price:   venue.Price,
		Games:   venue.Games,
		Rating:  venue.Rating,
		OwnerID: venue.OwnerID,
	}

	err := cs.storer.AddVenue(ctx, dbVenue)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error adding venue")
		return err
	}
	return nil
}

func (cs *UserOps) GetAllVenues(ctx context.Context) ([]domain.Venue, error) {
	venues, err := cs.storer.GetAllVenues(ctx)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting venues")
		return nil, err
	}
	responseVenues := make([]domain.Venue, len(venues))
	for i, venue := range venues {
		responseVenues[i] = domain.Venue{
			ID:      venue.ID,
			Name:    venue.Name,
			Address: venue.Address,
			City:    venue.City,
			State:   venue.State,
			Contact: venue.Contact,
			Email:   venue.Email,
			Opening: venue.Opening,
			Closing: venue.Closing,
			Price:   venue.Price,
			Games:   venue.Games,
			Rating:  venue.Rating,
			OwnerID: venue.OwnerID,
		}
	}
	return responseVenues, nil
}

func (cs *UserOps) GetVenue(ctx context.Context, id int) (domain.Venue, error) {
	fmt.Println("id", id)
	if id <= 0 {
		return domain.Venue{}, errors.New("invalid venue id")
	}
	venue, err := cs.storer.GetVenue(ctx, id)
	if err != nil && err == sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("no venue found")
		return domain.Venue{}, errors.New("no venue found")
	} else if err != nil {
		logger.WithField("err", err.Error()).Error("error getting venue")
		return domain.Venue{}, errors.New("error getting venue")
	}

	respVenue := domain.Venue{
		ID:      venue.ID,
		Name:    venue.Name,
		Address: venue.Address,
		City:    venue.City,
		State:   venue.State,
		Contact: venue.Contact,
		Email:   venue.Email,
		Opening: venue.Opening,
		Closing: venue.Closing,
		Price:   venue.Price,
		Games:   venue.Games,
		Rating:  venue.Rating,
		OwnerID: venue.OwnerID,
	}
	return respVenue, nil
}

func (cs *UserOps) UpdateVenue(ctx context.Context, venue domain.Venue, userID int, id int) error {
	dbVenue := db.Venue{
		ID:      venue.ID,
		Name:    venue.Name,
		Address: venue.Address,
		City:    venue.City,
		State:   venue.State,
		Contact: venue.Contact,
		Email:   venue.Email,
		Opening: venue.Opening,
		Closing: venue.Closing,
		Price:   venue.Price,
		Games:   venue.Games,
		Rating:  venue.Rating,
		OwnerID: venue.OwnerID,
	}

	err := cs.storer.UpdateVenue(ctx, dbVenue, userID, id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *UserOps) DeleteVenue(ctx context.Context, userID int, id int) error {
	err := cs.storer.DeleteVenue(ctx, userID, id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *UserOps) CheckAvailability(ctx context.Context, venueId int, date string) ([]domain.Slot, error) {
	fmt.Println("venueId", venueId, "date", date)
	slots, err := cs.storer.CheckAvailability(ctx, venueId, date)
	if err != nil {
		return nil, errors.New("error checking availability")
	}
	respSlots := make([]domain.Slot, len(slots))
	for i, slot := range slots {
		respSlots[i] = domain.Slot{
			VenueID:   slot.VenueID,
			Date:      slot.Date,
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
		}
	}

	return respSlots, nil
}

func (cs *UserOps) BookSlot(ctx context.Context, b domain.Booking) (float64, error) {
	dbBooking := db.Booking{
		ID:          b.ID,
		CustomerID:  b.CustomerID,
		VenueID:     b.VenueID,
		BookingDate: b.BookingDate,
		BookingTime: b.BookingTime,
		StartTime:   b.StartTime,
		EndTime:     b.EndTime,
		Game:        b.Game,
		AmountPaid:  b.AmountPaid,
	}

	price, err := cs.storer.BookSlot(ctx, dbBooking)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error booking slot")
		return 0.0, errors.New("error booking slot")
	}
	return price, nil

}

func (cs *UserOps) GetAllBookings(ctx context.Context, userId int) ([]domain.Booking, error) {
	bookings, err := cs.storer.GetAllBookings(ctx, userId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting bookings")
		return nil, errors.New("error getting bookings")
	}
	if len(bookings) == 0 {
		return nil, errors.New("no bookings found")
	}
	responseBookings := make([]domain.Booking, len(bookings))
	for i, booking := range bookings {
		responseBookings[i] = domain.Booking{
			ID:          booking.ID,
			CustomerID:  booking.CustomerID,
			VenueID:     booking.VenueID,
			BookingDate: booking.BookingDate,
			BookingTime: booking.BookingTime,
			StartTime:   booking.StartTime,
			EndTime:     booking.EndTime,
			Game:        booking.Game,
			AmountPaid:  booking.AmountPaid,
		}
	}
	return responseBookings, nil
}

func (cs *UserOps) GetBooking(ctx context.Context, id int) (domain.Booking, error) {
	booking, err := cs.storer.GetBooking(ctx, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booking")
		return domain.Booking{}, errors.New("error getting booking")
	}
	responseBooking := domain.Booking{
		ID:          booking.ID,
		CustomerID:  booking.CustomerID,
		VenueID:     booking.VenueID,
		BookingDate: booking.BookingDate,
		BookingTime: booking.BookingTime,
		StartTime:   booking.StartTime,
		EndTime:     booking.EndTime,
		Game:        booking.Game,
		AmountPaid:  booking.AmountPaid,
	}

	return responseBooking, err
}

func (cs *UserOps) CancelBooking(ctx context.Context, id int) error {
	err := cs.storer.CancelBooking(ctx, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking")
		return errors.New("error cancelling booking")
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
		return errors.New("user already exists")
	}
	return nil
}

func (cs *UserOps) CheckVenue(ctx context.Context, name string, contact string, email string) error {
	flag, err := cs.storer.CheckVenue(ctx, name, contact, email)
	if flag {
		logger.WithField("err", err.Error()).Error("error: venue already exists")
		return errors.New("venue already exists")
	}
	return nil
}
