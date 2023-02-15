package service

import (
	db "github.com/amancooks08/BookMySport/db"
)

type dependencies struct {
	CustomerServices Services
}

func InitDependencies() (deps dependencies, err error) {
	storer, err := db.Init()
	if err != nil {
		return
	}

	venueService := NewCustomerOps(storer)

	return dependencies{
		CustomerServices: venueService,
	}, err
}
