package db

import (
	"context"
)

type Storer interface {
	RegisterUser(context.Context, *User) error
	LoginUser(context.Context, string) (string, error)
	AddVenue(context.Context, *Venue) error
	GetAllVenues(context.Context) ([]*Venue, error)
	GetVenue(context.Context, string) (*Venue, error)
	UpdateVenue(context.Context, *Venue, int) error
	DeleteVenue(context.Context, int) error
	CheckAvailability(context.Context, int, string) ([]*Slot, error)
	BookSlot(context.Context, *Booking) error
	GetBooking(context.Context, int) (*Booking, error)
	GetAllBookings(context.Context, int) ([]*Booking, error)
	CancelBooking(context.Context, int) error
}
