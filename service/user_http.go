package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amancooks08/BookMySport/domain"
	logger "github.com/sirupsen/logrus"
)

type PingResponse struct {
	Message string `json:"message"`
}

func PingHandler(rw http.ResponseWriter, req *http.Request) {
	response := PingResponse{Message: "pong"}

	respBytes, err := json.Marshal(response)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error marshalling ping response")
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.Write(respBytes)
}

func RegisterCustomer(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		registerUser(rw, req, CustomerServices, "customer")
	})
}

func RegisterVenueOwner(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		registerUser(rw, req, CustomerServices, "venue_owner")
	})
}

func LoginUser(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var cu domain.UserLogin
		err := json.NewDecoder(req.Body).Decode(&cu)
		if err != nil {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		if cu.Email == "" || cu.Password == "" {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if validateEmail(cu.Email) {

			token, err := CustomerServices.LoginUser(req.Context(), cu.Email, cu.Password)
			if err != nil {
				msg := domain.Message{Message: fmt.Sprintf("%s", err)}
				respBytes, err := json.Marshal(msg)
				if err != nil {
					http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusUnauthorized)
				rw.Header().Add("Content-Type", "application/json")
				rw.Write(respBytes)
				return
			}
			if token != "" {
				response := domain.LoginResponse{
					Token:   token,
					Message: "Login Successful",
				}
				respBytes, err := json.Marshal(response)
				if err != nil {
					http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
					return
				}
				rw.Header().Add("Content-Type", "application/json")
				rw.Write(respBytes)
			} else {
				msg := domain.Message{Message: "error: invalid credentials"}
				respBytes, err := json.Marshal(msg)
				if err != nil {
					http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusUnauthorized)
				rw.Header().Add("Content-Type", "application/json")
				rw.Write(respBytes)
				return
			}

		} else {
			http.Error(rw, "invalid request payload", http.StatusBadRequest)
			return
		}
	})
}

// GetVenues function if called with id would return venue with that specific id and if not will return all venues

func GetVenues(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		venueID := GetVenueID(req)
		if venueID != 0 {
			venue, err := CustomerServices.GetVenue(req.Context(), venueID)
			if err != nil {
				msg := domain.Message{Message: fmt.Sprintf("%s", err.Error())}
				respBytes, err := json.Marshal(msg)
				if err != nil {
					http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
					return
				}
				// Add not found status
				rw.WriteHeader(http.StatusNotFound)
				rw.Header().Add("Content-Type", "application/json")
				rw.Write(respBytes)
				return
			}

			respBytes, err := json.Marshal(venue)
			if err != nil {
				http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
				return
			}

			rw.Header().Add("Content-Type", "application/json")
			rw.Write(respBytes)

		} else {
			venues, err := CustomerServices.GetAllVenues(req.Context())
			if err != nil {
				msg := domain.Message{Message: fmt.Sprintf("%s", err.Error())}
				respBytes, err := json.Marshal(msg)
				if err != nil {
					http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
					return
				}
				// Add not found status
				rw.WriteHeader(http.StatusNotFound)
				rw.Header().Add("Content-Type", "application/json")
				rw.Write(respBytes)
				return
			}

			if len(venues) == 0 {
				msg := domain.Message{Message: "no venues found"}
				respBytes, err := json.Marshal(msg)
				if err != nil {
					http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
					return
				}
				// Add not found status
				rw.Header().Add("Content-Type", "application/json")
				rw.WriteHeader(http.StatusNotFound)
				rw.Write(respBytes)
				return
			}
			respBytes, err := json.Marshal(venues)
			if err != nil {
				http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
				return
			}

			rw.Header().Add("Content-Type", "application/json")
			rw.Write(respBytes)
		}
	})
}

// Check availbility at a venue
func CheckAvailability(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		venueID := GetVenueID(req)

		// Check if date is present or not
		if req.URL.Query().Get("date") == "" {
			http.Error(rw, "please enter date if not entered or correct it if not added properly.", http.StatusBadRequest)
			return
		}
		date, err := time.Parse("2006-01-02", req.URL.Query().Get("date"))
		if err != nil {
			http.Error(rw, "invalid date format", http.StatusBadRequest)
			return
		}
		if date.After(time.Now().Truncate(24*time.Hour)) || date.Equal(time.Now().Truncate(24*time.Hour)) {
			availabileSlots, err := CustomerServices.CheckAvailability(req.Context(), venueID, date.Format("2006-01-02"))
			if err != nil && {
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
				return
			}

			respBytes, err := json.Marshal(availabileSlots)
			if err != nil {
				http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
				return
			}

			rw.Header().Add("Content-Type", "application/json")
			rw.Write(respBytes)
		} else {
			http.Error(rw, "invalid date - please selct an upcoming date.", http.StatusBadRequest)
			return
		}
	})
}
