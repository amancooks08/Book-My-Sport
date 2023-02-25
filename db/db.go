package db

import (
	"context"
)

type Storer interface {
	RegisterUser(context.Context, *User) error
	LoginUser(context.Context, string) (*LoginResponse, error)
	CheckUser(context.Context, string, string) (bool, error)
	AddVenue(context.Context, *Venue) error
	CheckVenue(context.Context, string, string, string) (bool, error)
	GetAllVenues(context.Context) ([]*Venue, error)
	GetVenue(context.Context, int) (*Venue, error)
	UpdateVenue(context.Context, *Venue, int) error
	DeleteVenue(context.Context, int) error
	CheckAvailability(context.Context, int, string) ([]*Slot, error)
	BookSlot(context.Context, *Booking) (float64, error)
	CalculatePrice(context.Context, int, int, string, string) (int, error)
	GetBooking(context.Context, int) (*Booking, error)
	GetAllBookings(context.Context, int) ([]*Booking, error)
	CancelBooking(context.Context, int) error
}

const (
	// User Queries
	RegisterUserQuery = `INSERT INTO "user" (name, contact, email, password, city, state, type)
    VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`
	LoginUserQuery = `SELECT id, password, type FROM "user" WHERE email = $1 `
	CheckUserQuery = `SELECT exists(SELECT 1 FROM "user" WHERE contact = $1 OR email = $2)`

	// Venue Queries
	InsertVenueQuery = `INSERT INTO "venue" (name, contact, city, state, address, email, opening_time, closing_time, price, rating)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`
	CheckVenueQuery      = `SELECT exists(SELECT 1 FROM "venue" WHERE name = $1 OR contact = $2 OR email = $3)`
	GetAllVenuesQuery    = `SELECT * FROM "venue"`
	GetVenueQuery        = `SELECT * FROM "venue" WHERE id = $1`
	UpdateVenueQuery     = `UPDATE "venue" SET name = $1, contact = $2, city = $3, state = $4, address = $5, opening_time = $6, closing_time = $7, price = $8, rating = $9 WHERE id = $10`
	DeleteVenueQuery     = `DELETE FROM "venue" WHERE id = $1`
	GetVenueTimingsQuery = `SELECT opening_time, closing_time FROM "venue" WHERE id = $1`

	// Slot Queries
	CheckAvailabilityQuery         = `"SELECT id, venue_id, start_time, end_time, date FROM "slots" WHERE venue_id = $1 date = $2 AND status = 'available'`
	InsertSlotQuery                = `INSERT INTO "slots" (venue_id, start_time, end_time, status, date) VALUES ($1, $2, $3, $4) RETURNING id`
	DeleteSlotQuery                = `DELETE FROM "slots" WHERE venue_id = $1`
	UpdateSlotStatusBookedQuery    = `UPDATE "slots" SET status = 'booked', booking_id = $1 WHERE start_time >= $2 AND end_time = $3 AND date = $4 AND venue_id = $5`
	UpdateSlotStatusAvailableQuery = `UPDATE "slots" SET status = 'available' WHERE booking_id = $1`
	UpdateSlotBookingQuery         = `UPDATE "slots" SET booking_id = NULL WHERE booking_id = $1`
	GetSlotsQuery                  = `SELECT id, venue_id, start_time, end_time, date FROM "slots" WHERE venue_id = $1 AND date = $2 AND status = 'available'`

	// Booking Queries
	SelectPriceQuery    = `SELECT price FROM "venue" WHERE id = $1`
	NumberOfSlotsQuery  = `SELECT COUNT(*) FROM "slots" WHERE booking_id = $1`
	CheckSlotStatusQuery = `SELECT exists(SELECT 1 FROM "slots" WHERE venue_id = $1 AND start_time >= $2 AND end_time <= $3 AND date = $4 AND status = 'available')`
	BookSlotQuery       = `INSERT INTO "booking" (booked_by, booked_at, booking_date, booking_time, start_time, end_time, game, amount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	GetBookingQuery     = `SELECT * FROM "slots" WHERE id = $1`
	GetAllBookingsQuery = `SELECT * FROM "booking" WHERE venue_id = $1`
	CancelBookingQuery  = `DELETE FROM "booking" WHERE id = $1`
	DeleteBookingQuery  = `DELETE FROM "booking" WHERE booked_at = $1`
)
