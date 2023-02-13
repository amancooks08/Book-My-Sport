package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Slot struct {
	Id        int    `json:"id"`
	VenueId   int    `json:"venue_id"`
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    bool 		`json:"status"`
}

func generateSlots(db *sqlx.DB, venueID int, startTime time.Time, endTime time.Time) error {
	day := time.Now().Format("2023-02-10")
	currentTime := startTime
	for currentTime.Before(endTime) {
		slotStart := day + " " + currentTime.Format("15:04:05")
		slotEnd := day + " " + currentTime.Add(1*time.Hour).Format("15:04:05")
		var slotExists bool
		err := db.QueryRow("SELECT exists(SELECT 1 FROM slots WHERE venue_id = $1 AND slot_start = $2 AND slot_end = $3)", venueID, slotStart, slotEnd).Scan(&slotExists)
		if err != nil {
			return err
		}
		zero := 0
		if !slotExists {
			_, err := db.Exec("INSERT INTO slots (venue_id, slot_start, slot_end, duration, status, date) VALUES ($1, $2, $3, $4, $5, $6)", venueID, slotStart, slotEnd, zero, day)
			if err != nil {
				return err
			}
		}
		currentTime = currentTime.Add(1 * time.Hour)
	}
	return nil
}
