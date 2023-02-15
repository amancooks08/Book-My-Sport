package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

/* The routing mechanism. Mux helps us define handler functions and the access methods */
func InitRouter(deps dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", PingHandler).Methods(http.MethodGet)

	api := router.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware)


	router.HandleFunc("/customer/register", RegisterCustomer(deps)).Methods(http.MethodPost)
	router.HandleFunc("/admin/register", RegisterAdmin(deps)).Methods(http.MethodPost)
	router.HandleFunc("/user/login", LoginUser(deps)).Methods(http.MethodPost)
	api.HandleFunc("/admin/venues/add", AddVenue(deps)).Methods(http.MethodPost)
	router.HandleFunc("/user/venues", GetAllVenues(deps)).Methods(http.MethodGet)
	router.HandleFunc("/user/venues/{venue_name}", GetVenue(deps)).Methods(http.MethodGet)
	api.HandleFunc("/admin/venues/{venue_id}", UpdateVenue(deps)).Methods(http.MethodPut)
	api.HandleFunc("/admin/venues/{venue_id}", DeleteVenue(deps)).Methods(http.MethodDelete)
	router.HandleFunc("/user/venues/{venue_id}/slots", CheckAvailability(deps)).Methods(http.MethodGet)
	api.HandleFunc("/customer/venues/{venue_id}/slots", BookSlot(deps)).Methods(http.MethodPost)
	api.HandleFunc("/customer/{user_id}/bookings", GetAllBookings(deps)).Methods(http.MethodGet)
	api.HandleFunc("/customer/bookings/{booking_id}", GetBooking(deps)).Methods(http.MethodGet)
	api.HandleFunc("/customer/bookings/{booking_id}/cancel", CancelBooking(deps)).Methods(http.MethodDelete)
	return
}
