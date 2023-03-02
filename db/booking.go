package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
)

type Booking struct {
	Id          int     `json:"id"`
	CustomerID  int     `json:"customer_id"`
	VenueID     int     `json:"venue_id"`
	BookingDate string  `json:"booking_date"`
	BookingTime string  `json:"booking_time"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Game        string  `json:"game"`
	AmountPaid  float64 `json:"amount"`
}

func (s *pgStore) BookSlot(ctx context.Context, b *Booking) (amount float64, err error) { //ad: named return can be used
	tx, err := s.db.Begin()
	if err != nil {
		logger.WithField("err", err.Error()).Error("error booking slot : Transaction Failed")
		return 0.0, ErrBeginTx
	}
	defer tx.Rollback()

	var price, slots int
	err = s.db.QueryRow(SelectPriceQuery, &b.VenueID).Scan(&price)
	if err != nil && err != sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error checking price")
		return 0.0, ErrCheckPrice
	}
	// Calculate the number of Slots byfinding the  duration
	st, _ := time.Parse("15:04:05", b.StartTime)
	et, _ := time.Parse("15:04:05", b.EndTime)
	slots = int(et.Sub(st).Hours())
	// Check if the game is present at the venue or not
	var game bool
	err = s.db.QueryRow(CheckGameQuery, b.VenueID, b.Game).Scan(&game)
	if err != nil && err != sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error checking game")
		return 0.0, ErrCheckGame
	}
	if !game {
		return 0.0, ErrGameNotAvailable
	}
	// Calculate Amount as a double value
	amount = float64(price * slots)
	var flag int
	err = s.db.QueryRow(CheckSlotStatusQuery, &b.VenueID, &b.StartTime, &b.EndTime, &b.BookingDate).Scan(&flag)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error checking slot status")
		return 0.0, ErrCheckSlotStatus
	}
	fmt.Println(flag)
	if flag == 0 {
		//Insert the booking details in the booking table and return id
		err = s.db.QueryRow(BookSlotQuery, &b.CustomerID, &b.VenueID, &b.BookingDate, &b.BookingTime, &b.StartTime, &b.EndTime, &b.Game, &amount).Scan(&b.Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error booking slot")
			return 0.0, ErrBookSlot
		}

		// Insert slots with booked status in the slot table
		generateSlots(s.db, b.VenueID, b.StartTime, b.EndTime, b.BookingDate, b.Id)
	} else {
		return 0.0, ErrSlotNotAvailable
	}
	tx.Commit()
	return amount, err

}

func (s *pgStore) GetBooking(ctx context.Context, id int) (booking *Booking, err error) {
	booking = &Booking{}
	err = s.db.QueryRow(GetBookingQuery, &id).Scan(&booking.Id, &booking.CustomerID, &booking.VenueID, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booking")
		return nil, ErrGetBooking
	}
	return booking, nil
}

func (s *pgStore) GetAllBookings(ctx context.Context, userId int) (bookings []*Booking, err error) {
	bookings = []*Booking{}

	rows, err := s.db.Query(GetAllBookingsQuery, &userId)
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("error: no bookings yet")
			return nil, ErrNoBookings
		}
		logger.WithField("err", err.Error()).Error("error getting all bookings")
		return nil, ErrGetBookings
	}

	for rows.Next() {
		booking := &Booking{}
		err := rows.Scan(&booking.Id, &booking.CustomerID, &booking.VenueID, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error getting all bookings")
			return nil, ErrGetBookings
		}

		booking = &Booking{
			Id:          booking.Id,
			CustomerID:  booking.CustomerID,
			VenueID:     booking.VenueID,
			BookingTime: booking.BookingTime[0:10] + " " + booking.BookingTime[11:19],
			BookingDate: booking.BookingDate[0:10],
			StartTime:   booking.StartTime[11:16],
			EndTime:     booking.EndTime[11:16],
			Game:        booking.Game,
			AmountPaid:  booking.AmountPaid,
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (s *pgStore) CancelBooking(ctx context.Context, id int) error {
	_, err := s.db.Exec(CancelBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking")
		return ErrCancelBooking
	}
	//Update the status of slots booked in the slots table to avialable

	_, err = s.db.Exec(UpdateSlotStatusAvailableQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating slot status")
		return ErrUpdateSlots
	}
	_, err = s.db.Exec(UpdateSlotBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating slot status to available")
		return ErrUpdateSlotsAvailable
	}
	return nil
}
