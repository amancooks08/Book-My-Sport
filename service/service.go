package service

import (
	"context"

	"github.com/amancooks08/BookMySport/db"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret@987")

type Services interface {
	RegisterUser(ctx context.Context, user *db.User) error
	LoginUser(ctx context.Context, email string, password string) (bool, error)
	AddVenue(ctx context.Context, venue *db.Venue) error
	GetAllVenues(ctx context.Context) ([]*db.Venue, error)
	GetVenue(ctx context.Context, name string) (*db.Venue, error)
	UpdateVenue(ctx context.Context, venue *db.Venue) error
	DeleteVenue(ctx context.Context, iname string) error
}

type UserOps struct {
	storer db.Storer
}

func NewCustomerOps(storer db.Storer) Services {
	return &UserOps{
		storer: storer,
	}
}


func (cs *UserOps) RegisterUser(ctx context.Context, user *db.User) error {
	user.Password, _ = HashPassword(user.Password)
	err := cs.storer.RegisterUser(ctx, user)
	return err
}

func (cs *UserOps) LoginUser(ctx context.Context, email string, password string) (bool, error) {
	flag, err := cs.storer.LoginUser(ctx, email)
	if bcrypt.CompareHashAndPassword([]byte(flag), []byte(password)) != nil {
		return false, err
	}
	return true, nil
}

func (cs *UserOps) AddVenue(ctx context.Context, venue *db.Venue) error {
	err := cs.storer.AddVenue(ctx, venue)
	return err
}

func (cs *UserOps) GetAllVenues(ctx context.Context) ([]*db.Venue, error) {
	venues, err := cs.storer.GetAllVenues(ctx)
	return venues, err
}

func (cs *UserOps) GetVenue(ctx context.Context, name string) (*db.Venue, error) {
	venue, err := cs.storer.GetVenue(ctx, name)
	return venue, err
}

func (cs *UserOps) UpdateVenue(ctx context.Context, venue *db.Venue) error {
	err := cs.storer.UpdateVenue(ctx, venue)
	return err
}

func (cs *UserOps) DeleteVenue(ctx context.Context, name string) error {
	err := cs.storer.DeleteVenue(ctx, name)
	return err
}

