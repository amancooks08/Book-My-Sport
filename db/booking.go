package db

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

type Booking struct {
	Id          int     `json:"id"`
	BookedBy    int     `json:"booked_by"`
	BookedAt    int     `json:"booked_at"`
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
	err = s.db.QueryRow(SelectPriceQuery, &b.BookedAt).Scan(&price)
	if err != nil && err != sql.ErrNoRows{
		logger.WithField("err", err.Error()).Error("error calculating price")
		return 0, err
	}

	// Calculate the number of Slots byfinding the  duration
	st, err := strconv.Atoi(b.StartTime[:2])
	if err != nil {
		logger.WithField("err", err.Error()).Error("error converting start time")
		return 0, errors.New("error converting start time")
	}
	et, err := strconv.Atoi(b.EndTime[:2])
	if err != nil {
		logger.WithField("err", err.Error()).Error("error converting end time")
		return 0, errors.New("error converting end time")
	}
	slots = et - st

	// Calculate Amount as a double value
	amount := float64(price * slots)
	var flag bool
	err = s.db.QueryRow(CheckSlotStatusQuery, &b.BookedAt, &b.StartTime, &b.EndTime, &b.BookingDate).Scan(&flag)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error checking slot status")
		return 0, errors.New("error checking slot status")
	}
	if flag {
		//Insert the booking details in the booking table
		_, err = s.db.Exec(BookSlotQuery, &b.BookedBy, &b.BookedAt, &b.BookingDate, &b.BookingTime, &b.StartTime, &b.EndTime, &b.Game, &amount)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error adding booking details")
			return 0.0, errors.New("error adding booking details")
		}

		// Update status of slots in "slots" table
		_, err = s.db.Exec(UpdateSlotStatusBookedQuery, &b.Id, &b.StartTime, &b.EndTime, &b.BookingDate, &b.BookedAt)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error updating slot status")
			return 0.0, errors.New("error updating slot status")
		}
	} else {
		return 0.0, errors.New("error: slot already booked")
	}	
	tx.Commit()
	return amount, err

}

func (s *pgStore) GetBooking(ctx context.Context, id int) (*Booking, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booking : Transaction Failed")
		return nil, err
	}
	defer tx.Rollback()

	booking := &Booking{}
	err = s.db.QueryRow(GetBookingQuery, &id).Scan(&booking.Id, &booking.BookedBy, &booking.BookedAt, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booking")
		return nil, err
	}
	tx.Commit()
	return booking, err
}

func (s *pgStore) GetAllBookings(ctx context.Context, userId int) ([]*Booking, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting all bookings : Transaction Failed")
		return nil, err
	}
	defer tx.Rollback()

	bookings := []*Booking{}

	rows, err := s.db.Query(GetAllBookingsQuery, &userId)

	// if err is norow error, return empty bookings

	if err != nil && err == sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error: no bookings yet")
		return nil, errors.New("error: no bookings yet")
	}
	tx.Commit()
	defer rows.Close()
	if rows == nil {
		return []*Booking{}, errors.New("error: no bookings yet")
	}
	for rows.Next() {
		booking := &Booking{}
		err := rows.Scan(&booking.Id, &booking.BookedBy, &booking.BookedAt, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error getting all bookings")
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, err
}

func (s *pgStore) CancelBooking(ctx context.Context, id int) error {
	_, err := s.db.Exec(CancelBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking")
		return err
	}
	//Update the status of slots booked in the slots table to avialable

	_, err = s.db.Exec(UpdateSlotStatusAvailableQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating slot status")
		return err
	}
	_, err = s.db.Exec(UpdateSlotBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating slot status to available")
		return err
	}
	return err
}
