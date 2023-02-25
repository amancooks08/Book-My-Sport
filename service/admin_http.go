package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/amancooks08/BookMySport/db"
	"github.com/gorilla/mux"
)

// Add a  venue
func AddVenue(CustomerServices Services) http.HandlerFunc {
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

		err = CustomerServices.CheckVenue(req.Context(), venue.Name, venue.Contact, venue.Email)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		if venue.Name == "" || venue.Address == "" || venue.City == "" || venue.State == "" {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if validateContact(venue.Contact) && validateEmail(venue.Email) {
			err := CustomerServices.AddVenue(req.Context(), &venue)
			if err != nil {
				http.Error(rw, fmt.Sprintf("%s1", err), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Write the response
		response := &Venue{
			Id:      venue.Id,
			Name:    venue.Name,
			Address: venue.Address,
			City:    venue.City,
			State:   venue.State,
			Contact: venue.Contact,
			Email:   venue.Email,
		}
		json_response, err := json.Marshal(response)
		if err != nil {
			http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(json_response)

	})
}

// Update a venue
func UpdateVenue(CustomerServices Services) http.HandlerFunc {
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
			err := CustomerServices.UpdateVenue(req.Context(), &venue, venueID)
			if err != nil {
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
	})
}

// Check availbility at a venue
func CheckAvailability(CustomerServices Services) http.HandlerFunc {
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
		// Check if date is present or not
		var date time.Time
		if req.URL.Query().Get("date") == "" {
			http.Error(rw, "Please enter date if not entered or correct it if not added properly.", http.StatusBadRequest)
			return
		}
		date, err = time.Parse("2006-01-02", req.URL.Query().Get("date"))
		if err != nil {
			http.Error(rw, "Invalid date format", http.StatusBadRequest)
			return
		}
		if date.After(time.Now().Truncate(24 * time.Hour)) || date.Equal(time.Now().Truncate(24 * time.Hour)) {
			availabileSlots, err := CustomerServices.CheckAvailability(req.Context(), venueID, date.Format("2006-01-02"))
			if err != nil {
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
				return
			}

			respBytes, err := json.Marshal(availabileSlots)
			if err != nil {
				http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
				return
			}

			rw.Header().Add("Content-Type", "application/json")
			rw.Write(respBytes)
		} else {
			http.Error(rw, "Invalid Date - Please selct an upcoming date.", http.StatusBadRequest)
			return
		}
	})
}

// Delete a venue
func DeleteVenue(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodDelete {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		vars := mux.Vars(req)
		venueID, err := strconv.Atoi(vars["venue_id"])
		//Check if "venue_id" key is not found in vars

		if err != nil {
			http.Error(rw, fmt.Sprint(err)+": Invalid ID", http.StatusBadRequest)
			return
		}
		err = CustomerServices.DeleteVenue(req.Context(), venueID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})
}
