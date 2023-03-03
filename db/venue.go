package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	logger "github.com/sirupsen/logrus"
)

type Venue struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	City    string   `json:"city"`
	State   string   `json:"state"`
	Contact string   `json:"contact"`
	Email   string   `json:"email"`
	Opening string   `json:"opening_time"`
	Closing string   `json:"closing_time"`
	Price   float64  `json:"price"`
	Games   []string `json:"games"`
	Rating  float64  `json:"rating"`
}

// AddVenue adds a venue to the database
func (s *pgStore) AddVenue(ctx context.Context, venue *Venue) error {
	err := s.db.QueryRow(InsertVenueQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating).Scan(&venue.ID)
	if err == errors.New(`pq: duplicate key value violates unique constraint \"venue_name_key\"`) {
		logger.WithField("err", err.Error()).Error("Error adding venue : Name already exists.")
		return ErrNameExists
	} else if err == errors.New(`pq: duplicate key value violates unique constraint \"venue_name_key\"`) {
		logger.WithField("err", err.Error()).Error("Error adding venue : Contact already exists.")
		return ErrContactExists
	} else if err == errors.New(`pq: duplicate key value violates unique constraint \"venue_email_key\"`) {
		logger.WithField("err", err.Error()).Error("Error adding venue : Email already exists.")
		return ErrEmailExists
	} else if err != nil {
		logger.WithField("err", err.Error()).Error("error adding venue")
		return ErrAddingVenue
	}
	return nil
}

// GetAllVenues returns all the venues in the database
func (s *pgStore) GetAllVenues(ctx context.Context) ([]*Venue, error) {
	venues := []*Venue{}
	rows, err := s.db.Query(GetAllVenuesQuery)
	if err != nil && err == sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error : no venues found")
		return nil, ErrNoVenues
	} else if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching all venues")
		return nil, ErrFetchingVenues
	}
	defer rows.Close()
	for rows.Next() {
		venue := &Venue{ID: 0, Name: "", Contact: "", City: "", State: "", Address: "", Email: "", Opening: "", Closing: "", Price: 0, Games: []string{}, Rating: 0}
		err = rows.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.City, &venue.State, &venue.Contact, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching all venues")
			return nil, ErrFetchingVenues
		}
		venue.Opening = venue.Opening[11:16]
		venue.Closing = venue.Closing[11:16]
		venues = append(venues, venue)
	}
	return venues, nil
}

// CheckVenue checks if a venue exists in the database

func (s *pgStore) CheckVenue(ctx context.Context, name string, contact string, email string) (bool, error) {
	var flag bool
	err := s.db.QueryRow(CheckVenueQuery, &name, &contact, &email).Scan(&flag)
	if err != nil && err == sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error checking venue")
		return false, ErrCheckVenue
	}
	return flag, ErrCheckVenue
}

// GetVenue returns a venue with the given id
func (s *pgStore) GetVenue(ctx context.Context, id int) (*Venue, error) {
	rows, err := s.db.Query(GetVenueQuery, &id)
	if err != nil && err == sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error: invalid venue id")
		return nil, ErrInvalidVID
	} else if err != nil {
		logger.WithField("err", err.Error()).Error("error fetching venue")
		return nil, ErrFetchingVenue
	}
	defer rows.Close()
	venue := &Venue{ID: 0, Name: "", Contact: "", City: "", State: "", Address: "", Email: "", Opening: "", Closing: "", Price: 0, Games: []string{}, Rating: 0}
	for rows.Next() {
		err = rows.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.City, &venue.State, &venue.Contact, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error fetching venue")
			return nil, ErrFetchingVenue
		}
	}
	return venue, nil
}

// UpdateVenue updates a venue in the database
func (s *pgStore) UpdateVenue(ctx context.Context, venue *Venue, id int) error {
	_, err := s.db.Exec(UpdateVenueQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating venue")
		return ErrUpdatingVenue
	}
	return nil
}

// DeleteVenue deletes a venue from the database
func (s *pgStore) DeleteVenue(ctx context.Context, id int) error {
	_, err := s.db.Exec(DeleteSlotQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return ErrDeletingVenue
	}

	_, err = s.db.Exec(DeleteBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return ErrDeletingVenue
	}

	_, err = s.db.Exec(DeleteVenueQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return ErrDeletingVenue
	}
	return nil
}

//Check availability of slots at a venue

func (s *pgStore) CheckAvailability(ctx context.Context, venueId int, date string) ([]*Slot, error) {
	var exists bool
	// Check if slots for the venue on that day exist
	err := s.db.QueryRow("SELECT exists(SELECT 1 FROM slots WHERE venue_id = $1 AND date = $2)", venueId, date).Scan(&exists)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error checking availability")
		return nil, ErrCheckAvailability
	}
	var venueOpen, venueClose string
	err = s.db.QueryRow(GetVenueTimingsQuery, venueId).Scan(&venueOpen, &venueClose)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting venue opening and closing times")
		return nil, ErrGetTimings
	}
	if !exists {
		// If slots don't exist, create them
		currentTime, err := time.Parse("15:04:05", venueOpen[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, ErrParseTime
		}
		venueClosing, err := time.Parse("15:04:05", venueClose[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, ErrParseTime
		}

		var slotList []*Slot
		for currentTime.Before(venueClosing) {
			slot := &Slot{VenueId: venueId, Date: date, StartTime: currentTime.Format("15:04"), EndTime: currentTime.Add(time.Hour).Format("15:04")}
			slotList = append(slotList, slot)
			currentTime = currentTime.Add(time.Hour)
		}
		return slotList, nil
	} else {
		// If slots existsquery for the venue on that day, return those which are available
		rows, err := s.db.Query(GetBookedSlotsQuery, &venueId, &date)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error getting booked slots")
			return nil, ErrBookedSlots
		}
		defer rows.Close()
		bookedSlots := []*Slot{}
		for rows.Next() {
			slot := &Slot{VenueId: 0, Date: "", StartTime: "", EndTime: ""}
			err = rows.Scan(&slot.VenueId, &slot.Date, &slot.StartTime, &slot.EndTime)
			slot.Date = slot.Date[0:10]
			slot.StartTime = slot.StartTime[11:16]
			slot.EndTime = slot.EndTime[11:16]
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error fetching booked slots")
				return nil, ErrBookedSlots
			}
			bookedSlots = append(bookedSlots, slot)
		}
		currentTime, err := time.Parse("15:04:05", venueOpen[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, ErrParseTime
		}
		venueClosing, err := time.Parse("15:04:05", venueClose[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, ErrParseTime
		}
		var slotList []*Slot
		for currentTime.Before(venueClosing) {
			slot := &Slot{VenueId: venueId, Date: date, StartTime: currentTime.Format("15:04"), EndTime: currentTime.Add(time.Hour).Format("15:04")}
			slotList = append(slotList, slot)
			currentTime = currentTime.Add(time.Hour)
		}
		for _, bookedSlot := range bookedSlots {
			for i, slot := range slotList {
				fmt.Println(slot.StartTime, bookedSlot.StartTime)
				if slot.StartTime == bookedSlot.StartTime {
					slotList = append(slotList[:i], slotList[i+1:]...)
				}
			}
		}
		return slotList, nil
	}
}
