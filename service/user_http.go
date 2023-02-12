package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

func RegisterCustomer(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		registerUser(rw, req, deps, "customer")
	})
}

func RegisterAdmin(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		registerUser(rw, req, deps, "admin")
	})
}

func LoginCustomer(deps dependencies) http.HandlerFunc {
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
			// Create a new JWT token
			token := jwt.New(jwt.SigningMethodHS256)

			// Set some claims
			claims := token.Claims.(jwt.MapClaims)
			claims["email"] = cu.Email
			claims["password"] = cu.Password
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

			// Sign the JWT token with the secret key
			tokenString, err := token.SignedString(secretKey)
			if err != nil {
				http.Error(rw, "Failed to sign JWT token", http.StatusInternalServerError)
				return
			}

			// Store the JWT token in a cookie
			http.SetCookie(rw, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: time.Now().Add(time.Hour * 24),
			})
			// Check if the user exists in the database
			if err != nil {
				http.Error(rw, "Failed to hash password", http.StatusInternalServerError)
				return
			}
			flag, err := deps.CustomerServices.LoginUser(req.Context(), cu.Email, cu.Password)
			if err != nil {
				http.Error(rw, fmt.Sprintf("%s", err), http.StatusUnauthorized)
				return
			}
			if flag {
				// Send a successful response
				rw.Write([]byte("User logged in successfully"))
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


// func bookVenue(deps Dependencies) http.HandlerFunc {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 		if req.Method != http.MethodPost {
// 			rw.WriteHeader(http.StatusMethodNotAllowed)
// 			return
// 		}

// 		// Decode the booking request from the request body
// 		var booking db.Booking
// 		err := json.NewDecoder(req.Body).Decode(&booking)
// 		if err != nil {
// 			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
// 			return
// 		}
// 		defer req.Body.Close()

// 		// Validate the booking request
// 		if booking.VenueID == 0 || booking.CustomerID == 0 || booking.StartTime.IsZero() || booking.EndTime.IsZero() {
// 			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
// 			return
// 		}

// 		// Check if the venue is available for booking
// 		available, err := deps.VenueServices.CheckAvailability(req.Context(), booking.VenueID, booking.StartTime, booking.EndTime)
// 		if err != nil {
// 			http.Error(rw, "Failed to check venue availability", http.StatusInternalServerError)
// 			return
// 		}
// 		if !available {
// 			http.Error(rw, "Venue is not available for booking", http.StatusBadRequest)
// 			return
// 		}

// 		// Store the booking in the database
// 		err = deps.CustomerService.BookVenue(req.Context(), booking)
// 		if err != nil {
// 			http.Error(rw, "Failed to book venue", http.StatusInternalServerError)
// 			return
// 		}

// 		// Send a successful response
// 		rw.Write([]byte("Venue booked successfully"))
// 	})
// }
