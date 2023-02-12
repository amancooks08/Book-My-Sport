package db

import (
	"context"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
)

type Venue struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Contact string    `json:"contact"`
	Email   string    `json:"email"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	Address string    `json:"address"`
	Opening time.Time `json:"opening_time"`
	Closing time.Time `json:"closing_time"`
	Price   float64   `json:"price"`
	Rating  float64   `json:"rating"`
}

func VenueExists(ctx context.Context, name string, s *pgStore) (bool, error) {
	sqlQuery := `SELECT EXISTS(SELECT 1 FROM "venue" WHERE name = $1)`
	var exists bool
	err := s.db.QueryRow(sqlQuery, name).Scan(&exists)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error checking for duplicate venue")
		return false, err
	}
	return exists, err
}

// AddVenue adds a venue to the database
func (s *pgStore) AddVenue(ctx context.Context, venue *Venue) error {

	// check for duplicate venue
	venueExists, err := VenueExists(ctx, venue.Name, s)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error checking for duplicate venue")
		return err
	}
	if venueExists {
		sqlQuery := `INSERT INTO "venue" (name, contact, city, state, address, opening_time, closing_time, price, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
		err := s.db.QueryRow(sqlQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating).Scan(&venue.Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error adding venue")
			return err
		}
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
		venue := &Venue{Id: 0, Name: "", Contact: "", City: "", State: "", Address: "", Opening: time.Time{}, Closing: time.Time{}, Price: 0, Rating: 0}
		err = rows.Scan(&venue.Id, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching venue")
			return nil, err
		}
		fmt.Printf(`Id : %v, Name : %v, Contact : %v, ` ,venue.Id, venue.Name, venue.Contact)
		venues = append(venues, venue)
	}
	return venues, err
}

// GetVenue returns a venue with the given id
func (s *pgStore) GetVenue(ctx context.Context, name string) (*Venue, error) {
	sqlQuery := `SELECT * FROM "venue" WHERE name LIKE "%$1%" OR name = $1`
	venue := &Venue{}
	err := s.db.QueryRow(sqlQuery, &name).Scan(&venue.Id, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting venue")
		return nil, err
	}
	return venue, err
}

// UpdateVenue updates a venue in the database
func (s *pgStore) UpdateVenue(ctx context.Context, venue *Venue) error {
	sqlQuery := `UPDATE "venue" SET name = $1, contact = $2, city = $3, state = $4, address = $5, opening_time = $6, closing_time = $7, price = $8, rating = $9 WHERE id = $10`
	_, err := s.db.Exec(sqlQuery, &venue.Name, &venue.Contact, &venue.City, &venue.State, &venue.Address, &venue.Opening, &venue.Closing, &venue.Price, &venue.Rating, &venue.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating venue")
		return err
	}
	return err
}

// DeleteVenue deletes a venue from the database
func (s *pgStore) DeleteVenue(ctx context.Context, name string) error {
	sqlQuery := `DELETE FROM "venue" WHERE name = $1`
	_, err := s.db.Exec(sqlQuery, &name)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting venue")
		return err
	}
	return err
}
