package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/amancooks08/BookMySport/db"
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

		var cu UserLogin
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
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusUnauthorized)
				return
			}
			if token != "" {
				response := &LoginResponse{
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
				msg := Message{Message: "error: invalid credentials"}
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

func GetAllVenues(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		venues, err := CustomerServices.GetAllVenues(req.Context())
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		if len(venues) == 0 {
			msg := Message{Message: "no venues found"}
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
		respBytes, err := json.Marshal(venues)
		if err != nil {
			http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}

func GetVenue(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		venueID, err := strconv.Atoi(req.URL.Query().Get("venue_id"))
		if err != nil || venueID <= 0 {
			msg := Message{Message: "invalid venue id"}
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
		venue, err := CustomerServices.GetVenue(req.Context(), venueID)
		if err != nil || err == db.ErrNoVenues {
			msg := Message{Message: "no venue found"}
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
		} else if err != nil {
			msg := Message{Message: err.Error()}
			respBytes, err := json.Marshal(msg)
			if err != nil {
				http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
				return
			}
			// Add not found status
			rw.WriteHeader(http.StatusInternalServerError)
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
	})

}
