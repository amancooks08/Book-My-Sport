package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
)

type Booking struct {
	ID          int     
	CustomerID  int     
	VenueID     int     
	BookingDate string  
	BookingTime string  
	StartTime   string  
	EndTime     string  
	Game        string  
	AmountPaid  float64 
}

func (s *pgStore) BookSlot(ctx context.Context, b Booking) (float64, error) {
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
	var gameExists bool
	err = s.db.QueryRow(CheckGameQuery, b.VenueID, b.Game).Scan(&gameExists)
	if err != nil && err != sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error checking game")
		return 0.0, ErrCheckGame
	}

	if !gameExists {
		return 0.0, ErrGameNotAvailable
	}

	// Calculate Amount as a double value
	amount := float64(price * slots)
	var flag int
	err = s.db.QueryRow(CheckSlotStatusQuery, &b.VenueID, &b.StartTime, &b.EndTime, &b.BookingDate).Scan(&flag)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error checking slot status")
		return 0.0, ErrCheckSlotStatus
	}

	fmt.Println(flag)
	if flag == 0 {
		//Insert the booking details in the booking table and return id
		err = s.db.QueryRow(BookSlotQuery, &b.CustomerID, &b.VenueID, &b.BookingDate, &b.BookingTime, &b.StartTime, &b.EndTime, &b.Game, &amount).Scan(&b.ID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error booking slot")
			return 0.0, ErrBookSlot
		}

		// Insert slots with booked status in the slot table
		generateSlots(s.db, b.VenueID, b.StartTime, b.EndTime, b.BookingDate, b.ID)
	} else {
		return 0.0, ErrSlotNotAvailable
	}
	
	tx.Commit()
	return amount, err

}

func (s *pgStore) GetBooking(ctx context.Context, id int) (Booking, error) {
	booking := Booking{}
	err := s.db.QueryRow(GetBookingQuery, &id).Scan(&booking.ID, &booking.CustomerID, &booking.VenueID, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booking")
		return Booking{}, ErrGetBooking
	}
	return booking, nil
}

func (s *pgStore) GetAllBookings(ctx context.Context, userId int) ([]Booking, error) {
	bookings := []Booking{}

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
		booking := Booking{}
		err := rows.Scan(&booking.ID, &booking.CustomerID, &booking.VenueID, &booking.BookingTime, &booking.BookingDate, &booking.StartTime, &booking.EndTime, &booking.Game, &booking.AmountPaid)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error getting all bookings")
			return nil, ErrGetBookings
		}

		booking = Booking{
			ID:          booking.ID,
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
	tx, err := s.db.Begin()
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking : Transaction Failed")
		return ErrBeginTx
	}
	defer tx.Rollback()

	// Remove booked slots from slots table
	_, err = s.db.Exec(RemoveBookedSlotsQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking")
		return ErrCancelBooking
	}

	// Remove booking from booking table
	_, err = s.db.Exec(CancelBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error cancelling booking")
		return ErrCancelBooking
	}
	tx.Commit()
	return nil
}
