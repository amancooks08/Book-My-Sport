package db

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

type Venue struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	City    string  `json:"city"`
	State   string  `json:"state"`
	Contact string  `json:"contact"`
	Email   string  `json:"email"`
	Opening string  `json:"opening_time"`
	Closing string  `json:"closing_time"`
	Price   float64 `json:"price"`
	Rating  float64 `json:"rating"`
}

// AddVenue adds a venue to the database
func (s *pgStore) AddVenue(ctx context.Context, venue *Venue) error {
	sqlQuery := `INSERT INTO venue (name, contact, city, state, address, email, opening_time, closing_time, price, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	err := s.db.QueryRow(sqlQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating).Scan(&venue.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error adding venue")
		return err
	}

	return err
}

// GetAllVenues returns all the venues in the database
func (s *pgStore) GetAllVenues(ctx context.Context) ([]*Venue, error) {
	sqlQuery := `SELECT * FROM "venue"`
	rows, err := s.db.Query(sqlQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting all venues")
		return nil, err
	}
	defer rows.Close()
	venues := []*Venue{}
	for rows.Next() {
		venue := &Venue{Id: 0, Name: "", Contact: "", City: "", State: "", Address: "", Email: "", Opening: "", Closing: "", Price: 0, Rating: 0}
		err = rows.Scan(&venue.Id, &venue.Name, &venue.Address, &venue.City, &venue.State, &venue.Contact, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching all venues")
			return nil, err
		}
		venues = append(venues, venue)
	}
	return venues, err
}

// GetVenue returns a venue with the given id
func (s *pgStore) GetVenue(ctx context.Context, name string) (*Venue, error) {
	sqlQuery := `SELECT * FROM "venue" WHERE name LIKE '%' || $1 || '%'`
	venue := &Venue{}
	rows, err := s.db.Query(sqlQuery, &name)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting venue")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&venue.Id, &venue.Name, &venue.Address, &venue.City, &venue.State, &venue.Contact, &venue.Email, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching venue")
			return nil, err
		}
	}

	return venue, err
}

// UpdateVenue updates a venue in the database
func (s *pgStore) UpdateVenue(ctx context.Context, venue *Venue, id int) error {
	sqlQuery := `UPDATE "venue" SET name = $1, contact = $2, city = $3, state = $4, address = $5, opening_time = $6, closing_time = $7, price = $8, rating = $9 WHERE id = $10`
	_, err := s.db.Exec(sqlQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating venue")
		return err
	}
	return err
}

// DeleteVenue deletes a venue from the database
func (s *pgStore) DeleteVenue(ctx context.Context, id int) error {
	sqlQuery := `DELETE FROM "venue" WHERE id = $1`
	_, err := s.db.Exec(sqlQuery, &id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return err
	}
	return err
}

//Check availability of slots at a venue

func (s *pgStore) CheckAvailability(ctx context.Context, venueId int, date string) ([]*Slot, error) {
	var exists bool
	// Check if slots for the venue on that day exist
	err := s.db.QueryRow("SELECT exists(SELECT 1 FROM slots WHERE venue_id = $1 AND date = $2)", venueId, date).Scan(&exists)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error checking availability")
		return nil, err
	}
	if !exists {
		// If slots don't exist, create them
		var venueOpen, venueClose string
		err = s.db.QueryRow("SELECT opening_time, closing_time FROM venue WHERE id = $1", venueId).Scan(&venueOpen, &venueClose)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error getting venue opening and closing times")
			return nil, err
		}
		err = generateSlots(s.db, venueId, venueOpen, venueClose, date)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error generating slots")
			return nil, err
		}
	}

	// Return all available slots for the venue on that day

	sqlQuery := `SELECT id, venue_id, start_time, end_time, date FROM "slots" WHERE venue_id = $1 AND date = $2 AND status = 'available'`
	rows, err := s.db.Query(sqlQuery, &venueId, &date)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting slots")
		return nil, err
	}
	defer rows.Close()
	slots := []*Slot{}
	for rows.Next() {
		slot := &Slot{Id: 0, VenueId: 0, Date: "", StartTime: "", EndTime: ""}
		err = rows.Scan(&slot.Id, &slot.VenueId, &slot.StartTime, &slot.EndTime, &slot.Date)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching slots")
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, err
}