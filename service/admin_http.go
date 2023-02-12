package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/amancooks08/BookMySport/db"
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
			err := deps.CustomerServices.AddVenue(req.Context(), &venue)
			if err != nil {
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
				return
			}
		}

		rw.WriteHeader(http.StatusCreated)
	})
}

func GetAllVenues(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		venues, err := deps.CustomerServices.GetAllVenues(req.Context())
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(venues)
		if err != nil {
			http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}

func GetVenue(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		venueID := req.URL.Query().Get("name")
		venue, err := deps.CustomerServices.GetVenue(req.Context(), venueID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(venue)
		if err != nil {
			http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
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

		if venue.Name == "" || venue.Address == "" || venue.City == "" || venue.State == "" {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if validateContact(venue.Contact) && validateEmail(venue.Email) {
			err := deps.CustomerServices.UpdateVenue(req.Context(), &venue)
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
		venueID := req.URL.Query().Get("name")
		err := deps.CustomerServices.DeleteVenue(req.Context(), venueID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})
}
