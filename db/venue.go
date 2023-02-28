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
	Id      int      `json:"id"`
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
	err := s.db.QueryRow(InsertVenueQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating).Scan(&venue.Id)
	if err == errors.New(`pq: duplicate key value violates unique constraint \"venue_name_key\"`) {
		logger.WithField("err", err.Error()).Error("Error adding venue : Name already exists.")
		return errors.New("error : name already exists")
	} else if err == errors.New(`pq: duplicate key value violates unique constraint \"venue_name_key\"`) {
		logger.WithField("err", err.Error()).Error("Error adding venue : Contact already exists.")
		return errors.New("error : contact already exists")
	} else if err == errors.New(`pq: duplicate key value violates unique constraint \"venue_email_key\"`) {
		logger.WithField("err", err.Error()).Error("Error adding venue : Email already exists.")
		return errors.New("error : email already exists")
	} else if err != nil {
		logger.WithField("err", err.Error()).Error("error adding venue")
		return errors.New("error adding venue")
	}
	return nil
}

// GetAllVenues returns all the venues in the database
func (s *pgStore) GetAllVenues(ctx context.Context) ([]*Venue, error) {
	rows, err := s.db.Query(GetAllVenuesQuery)
	if err != nil && err == sql.ErrNoRows {
		logger.WithField("err", err.Error()).Error("error : no venues found")
		return nil, errors.New("error : no venues found")
	} else if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching all venues")
		return nil, errors.New("error fetching all venues")
	}
	defer rows.Close()
	venues := []*Venue{}
	for rows.Next() {
		venue := &Venue{Id: 0, Name: "", Contact: "", City: "", State: "", Address: "", Email: "", Opening: "", Closing: "", Price: 0, Games: []string{}, Rating: 0}
		err = rows.Scan(&venue.Id, &venue.Name, &venue.Address, &venue.City, &venue.State, &venue.Contact, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching all venues")
			return nil, errors.New("error fetching all venues")
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
		return false, errors.New("error checking venue")
	}
	return flag, errors.New("error checking venue")
}

// GetVenue returns a venue with the given id
func (s *pgStore) GetVenue(ctx context.Context, id int) (*Venue, error) {
	venue := &Venue{}
	rows, err := s.db.Query(GetVenueQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting venue")
		return nil, errors.New("error getting venue")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&venue.Id, &venue.Name, &venue.Address, &venue.City, &venue.State, &venue.Contact, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching venue")
			return nil, errors.New("error fetching venue")
		}
	}
	return venue, nil
}

// UpdateVenue updates a venue in the database
func (s *pgStore) UpdateVenue(ctx context.Context, venue *Venue, id int) error {
	_, err := s.db.Exec(UpdateVenueQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Opening, &venue.Closing, &venue.Price, pq.Array(&venue.Games), &venue.Rating, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating venue")
		return errors.New("error updating venue")
	}
	return nil
}

// DeleteVenue deletes a venue from the database
func (s *pgStore) DeleteVenue(ctx context.Context, id int) error {
	_, err := s.db.Exec(DeleteSlotQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return errors.New("error deleting venue")
	}

	_, err = s.db.Exec(DeleteBookingQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return errors.New("error deleting venue")
	}

	_, err = s.db.Exec(DeleteVenueQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return errors.New("error deleting venue")
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
		return nil, errors.New("error checking availability")
	}
	var venueOpen, venueClose string
	err = s.db.QueryRow(GetVenueTimingsQuery, venueId).Scan(&venueOpen, &venueClose)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting venue opening and closing times")
		return nil, errors.New("error getting venue opening and closing times")
	}
	fmt.Println(exists)
	if !exists {
		// If slots don't exist, create them
		currentTime, err := time.Parse("15:04:05", venueOpen[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, errors.New("error parsing time")
		}
		venueClosing, err := time.Parse("15:04:05", venueClose[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, errors.New("error parsing time")
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
			return nil, errors.New("error getting booked slots")
		}
		defer rows.Close()
		bookedSlots := []*Slot{}
		for rows.Next() {
			slot := &Slot{VenueId: 0, Date: "", StartTime: "", EndTime: ""}
			err = rows.Scan(&slot.VenueId, &slot.StartTime, &slot.EndTime, &slot.Date)
			slot.Date = slot.Date[0:10]
			slot.StartTime = slot.StartTime[11:16]
			slot.EndTime = slot.EndTime[11:16]
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error fetching booked slots")
				return nil, errors.New("error fetching booked slots")
			}
			bookedSlots = append(bookedSlots, slot)
		}
		currentTime, err := time.Parse("15:04:05", venueOpen[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, errors.New("error parsing time")
		}
		venueClosing, err := time.Parse("15:04:05", venueClose[11:19])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing time")
			return nil, errors.New("error parsing time")
		}
		var slotList []*Slot
		for currentTime.Before(venueClosing) {
			slot := &Slot{VenueId: venueId, Date: date, StartTime: currentTime.Format("15:04"), EndTime: currentTime.Add(time.Hour).Format("15:04")}
			slotList = append(slotList, slot)
			currentTime = currentTime.Add(time.Hour)
		}
		for _, bookedSlot := range bookedSlots {
			for i, slot := range slotList {
				if slot.StartTime == bookedSlot.StartTime {
					slotList = append(slotList[:i], slotList[i+1:]...)
				}
			}
		}
		return slotList, nil
	}
}

// CalculatePrice calculates the price of a booking

func (s *pgStore) CalculatePrice(ctx context.Context, venueId int, bookingId int, startTime string, endTime string) (int, error) {
	var amount, price, slots int
	err := s.db.QueryRow(SelectPriceQuery, &venueId, &startTime, &endTime).Scan(&price)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error calculating price")
		return 0, errors.New("error calculating price")
	}

	err = s.db.QueryRow(NumberOfSlotsQuery, &bookingId).Scan(&slots)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error calculating price")
		return 0, errors.New("error calculating price")
	}
	
	amount = price * slots
	return amount, err
}
