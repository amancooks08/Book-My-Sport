package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
)

type Slot struct {
	Id        int    `json:"id"`
	VenueId   int    `json:"venue_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func generateSlots(db *sqlx.DB, venueID int, startTime string, endTime string, day string) error {
	currentTime, err := time.Parse("2006-01-02T15:04:05Z", startTime)
	if err != nil {
		return err
	}
	endingTime, err := time.Parse("2006-01-02T15:04:05Z", endTime)
	if err != nil {
		return err
	}

	for currentTime.Before(endingTime) {
		slotStart := currentTime.Format("15:04:05")
		slotEnd := currentTime.Add(1 * time.Hour).Format("15:04:05")
		var slotExists bool

		err := db.QueryRow(`SELECT exists(SELECT 1 FROM "slots" WHERE venue_id = $1 AND start_time = $2 AND end_time = $3 AND date = $4)`, venueID, slotStart, slotEnd, day).Scan(&slotExists)
		if err != nil {
			return err
		}

		availability := "available"

		if !slotExists {
			logger.Info("Generating slot")
			_, err := db.Exec("INSERT INTO slots (venue_id, start_time, end_time, status, date) VALUES ($1, $2, $3, $4, $5)", venueID, slotStart, slotEnd, availability, day)
			if err != nil {
				return err
			}
		}
		currentTime = currentTime.Add(1 * time.Hour)
	}
	return nil
}
