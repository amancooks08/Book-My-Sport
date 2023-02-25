package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func RegisterAdmin(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		registerUser(rw, req, CustomerServices, "admin")
	})
}

func LoginUser(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var cu customerLogin
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
					Token: token,
					Message: "Login Successful",
				}
				respBytes, err := json.Marshal(response)
				if err != nil {
					http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
					return
				}
				rw.Header().Add("Content-Type", "application/json")
				rw.Write(respBytes)
			} else {
				http.Error(rw, "Invalid credentials", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
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
		vars := mux.Vars(req)
		venueID, err := strconv.Atoi(vars["venue_id"])

		if err != nil {
			http.Error(rw, "Unable to parse VenueID", http.StatusBadRequest)
			return
		}

		venue, err := CustomerServices.GetVenue(req.Context(), venueID)
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
