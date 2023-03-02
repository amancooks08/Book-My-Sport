package db

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
)

type Slot struct {
	VenueId   int    `json:"venue_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func generateSlots(db *sqlx.DB, venueID int, startTime string, endTime string, day string, bookingId int) error {
	currentTime, _ := time.Parse("15:04:05", startTime)
	endingTime, _ := time.Parse("15:04:05", endTime)

	for currentTime.Before(endingTime) {
		slotStart := currentTime.Format("15:04:05")
		slotEnd := currentTime.Add(1 * time.Hour).Format("15:04:05")
		var slotExists bool

		err := db.QueryRow(`SELECT exists(SELECT 1 FROM "slots" WHERE venue_id = $1 AND start_time = $2 AND end_time = $3 AND date = $4)`, venueID, slotStart, slotEnd, day).Scan(&slotExists)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error checking if slot exists")
			return errors.New("error checking if slot exists")
		}

		const availability = "booked"
		if !slotExists {
			_, err := db.Exec("INSERT INTO slots (venue_id, start_time, end_time, status, date, booking_id) VALUES ($1, $2, $3, $4, $5, $6)", venueID, slotStart, slotEnd, availability, day, bookingId)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error inserting slot")
				return errors.New("error inserting slot")
			}
		}
		currentTime = currentTime.Add(1 * time.Hour)
	}
	return nil
}
