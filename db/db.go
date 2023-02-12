package db

import (
	"context"
	"time"
)

type Storer interface {
	RegisterUser(context.Context, *User) error
	LoginUser(context.Context, string) (string, error)
	AddVenue(context.Context, *Venue) error
	GetAllVenues(context.Context) ([]*Venue, error)
	GetVenue(context.Context, string) (*Venue, error)
	UpdateVenue(context.Context, *Venue, int) error
	DeleteVenue(context.Context, int) error
	CheckAvailability(context.Context, int, time.Time) ([]*Slot, error)
	// BookVenue(context.Context, *Booking) error
	// GetBookings(context.Context, string) ([]Booking, error)
	// GetAllBookings(context.Context) ([]Booking, error)
	// CancelBooking(context.Context, string) error

}
