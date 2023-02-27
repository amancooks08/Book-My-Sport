package server

import (
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/amancooks08/BookMySport/service"
)

/* The routing mechanism. Mux helps us define handler functions and the access methods */
func InitRouter(deps dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", service.PingHandler).Methods(http.MethodGet)

	router.HandleFunc("/customer/register", service.RegisterCustomer(deps.CustomerServices)).Methods(http.MethodPost)
	router.HandleFunc("/venue_owner/register", service.RegisterVenueOwner(deps.CustomerServices)).Methods(http.MethodPost)
	router.HandleFunc("/user/login", service.LoginUser(deps.CustomerServices)).Methods(http.MethodPost)
	router.HandleFunc("/venue_owner/venues/add", authMiddleware(service.AddVenue(deps.CustomerServices))).Methods(http.MethodPost)
	router.HandleFunc("/user/venues", service.GetAllVenues(deps.CustomerServices)).Methods(http.MethodGet)
	router.HandleFunc("/user/venues/{venue_id}", service.GetVenue(deps.CustomerServices)).Methods(http.MethodGet)
	router.HandleFunc("/venue_owner/venues/{venue_id}", authMiddleware(service.UpdateVenue(deps.CustomerServices))).Methods(http.MethodPut)
	router.HandleFunc("/venue_owner/venues/{venue_id}", authMiddleware(service.DeleteVenue(deps.CustomerServices))).Methods(http.MethodDelete)
	router.HandleFunc("/user/venues/{venue_id}/slots", service.CheckAvailability(deps.CustomerServices)).Methods(http.MethodGet)
	router.HandleFunc("/customer/venues/{venue_id}/book", authMiddleware(service.BookSlot(deps.CustomerServices))).Methods(http.MethodPost)
	router.HandleFunc("/customer/bookings", authMiddleware(service.GetAllBookings(deps.CustomerServices))).Methods(http.MethodGet)
	router.HandleFunc("/customer/bookings/{booking_id}", authMiddleware(service.GetBooking(deps.CustomerServices))).Methods(http.MethodGet)
	router.HandleFunc("/customer/bookings/{booking_id}/cancel", authMiddleware(service.CancelBooking(deps.CustomerServices))).Methods(http.MethodDelete)
	return
}
