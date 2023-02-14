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

func BookSlot(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var booking db.Booking
		err := json.NewDecoder(req.Body).Decode(&booking)
		if err != nil {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		booking.BookingTime = time.Now().Format("2006-01-02 15:04:05.999999-07")
		if booking.BookedBy == 0 || booking.BookedAt == 0 || booking.BookingDate < time.Now().Format("2006-01-02") || booking.BookingTime == "" {
			http.Error(rw, "Error : Please enter correct details.", http.StatusBadRequest)
			return
		}

		err = deps.CustomerServices.BookSlot(req.Context(), &booking)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		type jsonResponse struct {
			Reponse string
		}

		response := jsonResponse{Reponse: "Booking successful."}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)
		rw.WriteHeader(http.StatusCreated)
	})
}

// Get all bookings for a user

func GetAllBookings(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(req)
		userID, err := strconv.Atoi(vars["user_id"])
		if err != nil {
			http.Error(rw, "Invalid user ID", http.StatusBadRequest)
			return
		}

		bookings, err := deps.CustomerServices.GetAllBookings(req.Context(), userID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		type jsonResponse struct {
			Bookings []*db.Booking
		}

		response := jsonResponse{Bookings: bookings}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)
		rw.WriteHeader(http.StatusOK)
	})
}

// Get a Specific Booking

func GetBooking(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(req)
		bookingID, err := strconv.Atoi(vars["booking_id"])
		if err != nil {
			http.Error(rw, "Invalid booking ID", http.StatusBadRequest)
			return
		}

		booking, err := deps.CustomerServices.GetBooking(req.Context(), bookingID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		type jsonResponse struct {
			Booking *db.Booking
		}

		response := jsonResponse{Booking: booking}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)
		rw.WriteHeader(http.StatusOK)
	})
}

//Cancel Existing Booking

func CancelBooking(deps dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(req)
		bookingID, err := strconv.Atoi(vars["booking_id"])
		if err != nil {
			http.Error(rw, "Invalid booking ID", http.StatusBadRequest)
			return
		}

		err = deps.CustomerServices.CancelBooking(req.Context(), bookingID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		type jsonResponse struct {
			Reponse string
		}

		response := jsonResponse{Reponse: fmt.Sprintf("Booking cancelled successfully.")}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)
		rw.WriteHeader(http.StatusOK)
	})
}
