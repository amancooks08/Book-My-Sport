package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/amancooks08/BookMySport/db"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

func AddVenue(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var venue db.Venue
		err := json.NewDecoder(req.Body).Decode(&venue)
		if err != nil {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		if venue.Name == "" || venue.Address == "" || venue.City == "" || venue.State == "" {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if validateContact(venue.Contact) && validateEmail(venue.Email) {
			logger.Info("Adding venue: ", venue.Name)
			err := deps.CustomerServices.AddVenue(req.Context(), &venue)
			if err != nil {
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	})
}

func UpdateVenue(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var venue db.Venue
		err := json.NewDecoder(req.Body).Decode(&venue)
		if err != nil {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		vars := mux.Vars(req)
		venueID, err := strconv.Atoi(vars["venue_id"])
		if err != nil {
			http.Error(rw, fmt.Sprint(err)+": Invalid ID", http.StatusBadRequest)
			return
		}
		if venue.Name == "" || venue.Address == "" || venue.City == "" || venue.State == "" {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if validateContact(venue.Contact) && validateEmail(venue.Email) {
			err := deps.CustomerServices.UpdateVenue(req.Context(), &venue, venueID)
			if err != nil {
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
	})
}

func DeleteVenue(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		vars := mux.Vars(req)
		venueID, err := strconv.Atoi(vars["venue_id"])
		if err != nil {
			http.Error(rw, fmt.Sprint(err)+": Invalid ID", http.StatusBadRequest)
			return
		}
		err = deps.CustomerServices.DeleteVenue(req.Context(), venueID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})
}

// Check availbility at a venue

func CheckAvailability(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		vars := mux.Vars(req)
		venueID, err := strconv.Atoi(vars["venue_id"])
		if err != nil {
			http.Error(rw, fmt.Sprint(err)+": Invalid ID", http.StatusBadRequest)
			return
		}

		date, err := time.Parse("2023-02-10", req.URL.Query().Get("date"))
		if err != nil {
			http.Error(rw, fmt.Sprint(err)+": Invalid date", http.StatusBadRequest)
			return
		}
		availability, err := deps.CustomerServices.CheckAvailability(req.Context(), venueID, date)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(availability)
		if err != nil {
			http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}
