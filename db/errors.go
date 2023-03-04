package db

import (
	"errors"
)

// Custom Errors
var (
	ErrNameExists    = errors.New("name already exists")
	ErrContactExists = errors.New("contact already exists")
	ErrEmailExists   = errors.New("email already exists")

	ErrAddingVenue    = errors.New("error adding venue")
	ErrNoVenues       = errors.New("no venues found")
	ErrNoVenue        = errors.New("no venue found")
	ErrFetchingVenues = errors.New("error fetching all venues")
	ErrCheckVenue     = errors.New("error checking venue")
	ErrInvalidVID     = errors.New("invalid venue id")
	ErrFetchingVenue  = errors.New("error fetching venue")
	ErrUpdatingVenue  = errors.New("error updating venue")
	ErrDeletingVenue  = errors.New("error deleting venue")
	ErrVenueOwnerNotFound = errors.New("user is not the owner of this venue")
	ErrCheckVenueOwner = errors.New("error checking venue owner")

	ErrCheckAvailability = errors.New("error checking availability")
	ErrGetTimings        = errors.New("error getting venue opening and closing times")
	ErrBookedSlots       = errors.New("error getting booked slots")
	ErrCalculatePrice    = errors.New("error calculating price")
	ErrParseTime         = errors.New("error parsing time")

	ErrUserNotExists = errors.New("no user with that email id exist")
	ErrLogin         = errors.New("error logging in")

	ErrCheckUser    = errors.New("error checking user")
	ErrRegisterUser = errors.New("error registering user")

	ErrBeginTx          = errors.New("error beginning transaction")
	ErrCommitTx         = errors.New("error committing transaction")
	ErrCheckPrice       = errors.New("error checking price")
	ErrCheckGame        = errors.New("error checking game")
	ErrGameNotAvailable = errors.New("game not available at this venue")
	ErrCheckSlotStatus  = errors.New("error checking slot status")
	ErrBookSlot         = errors.New("error booking slot")
	ErrSlotNotAvailable = errors.New("slot already booked")

	ErrGetBooking           = errors.New("error getting booking")
	ErrNoBookings           = errors.New("no bookings found")
	ErrGetBookings          = errors.New("error getting all bookings")
	ErrCancelBooking        = errors.New("error cancelling booking")
	ErrUpdateSlots          = errors.New("error updating slots")
	ErrUpdateSlotsAvailable = errors.New("error updating slots status to available")
)
