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

func BookSlot(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Get The userId from the JWT token
		var booking db.Booking
		err := json.NewDecoder(req.Body).Decode(&booking)
		if err != nil {
			http.Error(rw, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		booking.BookedBy, booking.BookedAt = GetUserVenueId(req, rw)
		booking.BookingTime = time.Now().Format("2006-01-02 15:04:05.999999-07")
		if booking.BookingDate <= time.Now().Format("2006-01-02") || booking.BookingTime == "" || booking.StartTime == "" || booking.EndTime == "" || booking.Game == "" {
			http.Error(rw, "Error : Please enter correct details.", http.StatusBadRequest)
			return
		}

		amount, err := CustomerServices.BookSlot(req.Context(), &booking)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		type jsonResponse struct {
			Reponse string
			Amount  float64
		}

		response := jsonResponse{Reponse: "Booking successful.", Amount: amount}
		// rw.WriteHeader(http.StatusCreated)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)

	})
}

// Get all bookings for a user

func GetAllBookings(CustomerServices Services) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userID, _ := GetUserVenueId(req, rw)

		bookings, err := CustomerServices.GetAllBookings(req.Context(), userID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		
		// If bookings is empty
		if len(bookings) == 0 {
			http.Error(rw, "No bookings found", http.StatusNotFound)
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

func GetBooking(CustomerServices Services) http.HandlerFunc {
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

		booking, err := CustomerServices.GetBooking(req.Context(), bookingID)
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

func CancelBooking(CustomerServices Services) http.HandlerFunc {
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

		err = CustomerServices.CancelBooking(req.Context(), bookingID)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		type jsonResponse struct {
			Reponse string
		}

		response := jsonResponse{Reponse: "Booking cancelled successfully."}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)
		rw.WriteHeader(http.StatusOK)
	})
}
