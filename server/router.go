package server

import (
	"net/http"

	service "github.com/amancooks08/BookMySport/service"
	"github.com/gorilla/mux"
)

/* The routing mechanism. Mux helps us define handler functions and the access methods */
func InitRouter(deps dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", service.PingHandler).Methods(http.MethodGet)

	router.HandleFunc("/customer/register", service.RegisterCustomer(deps.CustomerServices)).Methods(http.MethodPost)
	router.HandleFunc("/venue_owner/register", service.RegisterVenueOwner(deps.CustomerServices)).Methods(http.MethodPost)
	router.HandleFunc("/user/login", service.LoginUser(deps.CustomerServices)).Methods(http.MethodPost)
	router.HandleFunc("/venue_owner/venues", authMiddleware(service.AddVenue(deps.CustomerServices))).Methods(http.MethodPost)
	router.HandleFunc("/user/venues", authMiddleware(service.GetVenues(deps.CustomerServices))).Methods(http.MethodGet)
	router.HandleFunc("/venue_owner/venues", authMiddleware(service.UpdateVenue(deps.CustomerServices))).Methods(http.MethodPut)
	router.HandleFunc("/venue_owner/venues", authMiddleware(service.DeleteVenue(deps.CustomerServices))).Methods(http.MethodDelete)
	router.HandleFunc("/user/venues/slots", authMiddleware(service.CheckAvailability(deps.CustomerServices))).Methods(http.MethodGet)
	router.HandleFunc("/customer/venues/book", authMiddleware(service.BookSlot(deps.CustomerServices))).Methods(http.MethodPost)
	router.HandleFunc("/customer/bookings", authMiddleware(service.GetAllBookings(deps.CustomerServices))).Methods(http.MethodGet)
	router.HandleFunc("/customer/bookings/{bookingID}", authMiddleware(service.GetBooking(deps.CustomerServices))).Methods(http.MethodGet)
	router.HandleFunc("/customer/bookings/{bookingID}/cancel", authMiddleware(service.CancelBooking(deps.CustomerServices))).Methods(http.MethodDelete)
	return
}
