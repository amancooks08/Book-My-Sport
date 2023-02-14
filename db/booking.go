package db

import (
	"context"

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

func (s *pgStore) BookSlot(ctx context.Context, b *Booking) error {
	sqlQuery := "INSERT INTO booking (booked_by, booked_at, booking_date, booking_time, start_time, end_time, game, amount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := s.db.Exec(sqlQuery, &b.BookedBy, &b.BookedAt, &b.BookingDate, &b.BookingTime, &b.StartTime, &b.EndTime, &b.Game, &b.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error booking slot")
		return err
	}
	//Update the status of slots booked in the slots table

	sqlQuery = "UPDATE slots SET status = 'booked' WHERE start_time >= $1 AND end_time <= $2"
	rows, err := s.db.Query(sqlQuery, &b.StartTime, &b.EndTime)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating slot status to booked")
		return err
	}
	defer rows.Close()
	return err
}

func (s *pgStore) GetBooking(ctx context.Context, id int) (*Booking, error) {
	sqlQuery := "SELECT * FROM booking WHERE id = $1"
	booking := &Booking{}
	err := s.db.QueryRow(sqlQuery, &id).Scan(&booking.Id, &booking.BookedBy, &booking.BookedAt, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting booking")
		return nil, err
	}
	return booking, err
}

func (s *pgStore) GetAllBookings(ctx context.Context, userId int) ([]*Booking, error) {
	sqlQuery := "SELECT * FROM booking WHERE booked_by = $1"
	bookings := []*Booking{}
	rows, err := s.db.Query(sqlQuery, &userId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting all bookings")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		booking := &Booking{}
		err := rows.Scan(&booking.Id, &booking.BookedBy, &booking.BookedAt, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error getting all bookings")
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, err
}

func (s *pgStore) CancelBooking(ctx context.Context, id int) error {
	sqlQuery := "DELETE FROM booking WHERE id = $1"
	_, err := s.db.Exec(sqlQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error cancelling booking")
		return err
	}
	//Update the status of slots booked in the slots table to avialable

	sqlQuery = "UPDATE slots SET status = 'available' WHERE booking_id = $1"
	rows, err := s.db.Query(sqlQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating slot status")
		return err
	}
	sqlQuery = "UPDATE slots SET booking_id = NULL WHERE booking_id = $1"
	rows, err = s.db.Query(sqlQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating slot status to available")
		return err
	}
	defer rows.Close()
	return err
}