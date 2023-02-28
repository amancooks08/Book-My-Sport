package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
)

type Booking struct {
	Id          int     `json:"id"`
	CustomerId  int     `json:"customer_id"`
	VenueId     int     `json:"venue_id"`
	BookingDate string  `json:"booking_date"`
	BookingTime string  `json:"booking_time"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Game        string  `json:"game"`
	AmountPaid  float64 `json:"amount"`
}

func (s *pgStore) BookSlot(ctx context.Context, b *Booking) (float64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logger.WithField("err", err.Error()).Error("error booking slot : Transaction Failed")
		return 0.0, err
	}
	defer tx.Rollback()

	var price, slots int
	err = s.db.QueryRow(SelectPriceQuery, &b.VenueId).Scan(&price)
	if err != nil && err != sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error calculating price")
		return 0, err
	}
	// Calculate the number of Slots byfinding the  duration
	st, _ := time.Parse("15:04:00", b.StartTime)
	et, _ := time.Parse("15:04:00", b.EndTime)
	slots = int(et.Sub(st).Hours())
	// Check if the game is present at the venue or not
	var game bool
	err = s.db.QueryRow(CheckGameQuery, b.VenueId, b.Game).Scan(&game)
	if err != nil && err != sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error checking game")
		return 0, errors.New("error checking game")
	}
	if !game {
		return 0, errors.New("error: game not available at this venue")
	}
	// Calculate Amount as a double value
	amount := float64(price * slots)
	var flag int
	err = s.db.QueryRow(CheckSlotStatusQuery, &b.VenueId, &b.StartTime, &b.EndTime, &b.BookingDate).Scan(&flag)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error checking slot status")
		return 0, errors.New("error checking slot status")
	}
	fmt.Println(flag)
	if flag == 0 {
		//Insert the booking details in the booking table and return id
		err = s.db.QueryRow(BookSlotQuery, &b.CustomerId, &b.VenueId, &b.BookingDate, &b.BookingTime, &b.StartTime, &b.EndTime, &b.Game, &amount).Scan(&b.Id)
		if err != nil {

			logger.WithField("err", err.Error()).Error("error booking slot")
			return 0.0, errors.New("error booking slot")
		}

		// Insert slots with booked status in the slot table
		generateSlots(s.db, b.VenueId, b.StartTime, b.EndTime, b.BookingDate, b.Id)
	} else {
		return 0.0, errors.New("error: slot already booked")
	}
	tx.Commit()
	return amount, err

}

func (s *pgStore) GetBooking(ctx context.Context, id int) (*Booking, error) {
	booking := &Booking{}
	err := s.db.QueryRow(GetBookingQuery, &id).Scan(&booking.Id, &booking.CustomerId, &booking.VenueId, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booking")
		return nil, errors.New("error getting booking")
	}
	return booking, nil
}

func (s *pgStore) GetAllBookings(ctx context.Context, userId int) ([]*Booking, error) {
	bookings := []*Booking{}

	rows, err := s.db.Query(GetAllBookingsQuery, &userId)
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("error: no bookings yet")
			return nil, errors.New("error: no bookings yet")
		}
		logger.WithField("err", err.Error()).Error("error getting all bookings")
		return nil, errors.New("error getting all bookings")
	}

	for rows.Next() {
		booking := &Booking{}
		err := rows.Scan(&booking.Id, &booking.CustomerId, &booking.VenueId, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error getting all bookings")
			return nil, errors.New("error getting all bookings")
		}
		booking.BookingDate = booking.BookingDate[0:10]
		booking.BookingTime = booking.BookingTime[0:10] + " " + booking.BookingTime[11:19]
		booking.StartTime = booking.StartTime[11:16]
		booking.EndTime = booking.EndTime[11:16]
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (s *pgStore) CancelBooking(ctx context.Context, id int) error {
	_, err := s.db.Exec(CancelBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking")
		return errors.New("error cancelling booking")
	}
	//Update the status of slots booked in the slots table to avialable

	_, err = s.db.Exec(UpdateSlotStatusAvailableQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating slot status")
		return errors.New("error updating slot status")
	}
	_, err = s.db.Exec(UpdateSlotBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating slot status to available")
		return errors.New("error updating slot status to available")
	}
	return nil
}
