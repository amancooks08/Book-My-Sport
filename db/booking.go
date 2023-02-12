package db

import (
	"context"
	"time"

	logger "github.com/sirupsen/logrus"
)

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

func (s *pgStore) BookVenue(ctx context.Context, b *Booking) error {
	sqlQuery := "INSERT INTO bookings (booked_by, booked_at, booking_time, start_time, end_time, game, amount) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := s.db.Exec(sqlQuery, &b.BookedBy, &b.BookedAt, &b.Time, &b.StartTime, &b.EndTime, &b.Game, &b.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error booking venue")
		return err
	}
	return err
}

func (s *pgStore) GetBooking(ctx context.Context, id int) (*Booking, error) {
	sqlQuery := "SELECT * FROM bookings WHERE id = $1"
	booking := &Booking{}
	err := s.db.QueryRow(sqlQuery, &id).Scan(&booking.Id, &booking.BookedBy, &booking.BookedAt, &booking.Time, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting booking")
		return nil, err
	}
	return booking, err
}
