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

	// Version 1 API management
	// v1 := fmt.Sprintf("application/vnd", config.AppName())
	router.HandleFunc("/customer/register", RegisterCustomer(deps)).Methods(http.MethodPost)
	router.HandleFunc("/admin/register", RegisterAdmin(deps)).Methods(http.MethodPost)
	router.HandleFunc("/user/login", LoginCustomer(deps)).Methods(http.MethodPost)
	router.HandleFunc("/admin/venues/add", AddVenue(deps)).Methods(http.MethodPost)
	router.HandleFunc("/user/venues", GetAllVenues(deps)).Methods(http.MethodGet)
	router.HandleFunc("/user/venues/{name}", GetVenue(deps)).Methods(http.MethodGet)
	router.HandleFunc("/admin/venues/{name}", UpdateVenue(deps)).Methods(http.MethodPut)
	router.HandleFunc("/admin/venues/{name}", DeleteVenue(deps)).Methods(http.MethodDelete)
	return
}
