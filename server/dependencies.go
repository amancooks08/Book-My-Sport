package server

import (
	db "github.com/amancooks08/BookMySport/db"
	service "github.com/amancooks08/BookMySport/service"
)

type dependencies struct {
	CustomerServices service.Services
}

func InitDependencies() (deps dependencies, err error) {
	storer, err := db.Init()
	if err != nil {
		return
	}

	venueService := service.NewCustomerOps(storer)

	return dependencies{
		CustomerServices: venueService,
	}, err
}
