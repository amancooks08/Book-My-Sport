// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	db "github.com/amancooks08/BookMySport/db"
	mock "github.com/stretchr/testify/mock"
)

// Services is an autogenerated mock type for the Services type
type Services struct {
	mock.Mock
}

// AddVenue provides a mock function with given fields: ctx, venue
func (_m *Services) AddVenue(ctx context.Context, venue *db.Venue) error {
	ret := _m.Called(ctx, venue)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *db.Venue) error); ok {
		r0 = rf(ctx, venue)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BookSlot provides a mock function with given fields: ctx, b
func (_m *Services) BookSlot(ctx context.Context, b *db.Booking) (float64, error) {
	ret := _m.Called(ctx, b)

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *db.Booking) (float64, error)); ok {
		return rf(ctx, b)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *db.Booking) float64); ok {
		r0 = rf(ctx, b)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *db.Booking) error); ok {
		r1 = rf(ctx, b)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CancelBooking provides a mock function with given fields: ctx, id
func (_m *Services) CancelBooking(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckAvailability provides a mock function with given fields: ctx, id, date
func (_m *Services) CheckAvailability(ctx context.Context, id int, date string) ([]*db.Slot, error) {
	ret := _m.Called(ctx, id, date)

	var r0 []*db.Slot
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) ([]*db.Slot, error)); ok {
		return rf(ctx, id, date)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string) []*db.Slot); ok {
		r0 = rf(ctx, id, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*db.Slot)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, id, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckUser provides a mock function with given fields: ctx, email, contact
func (_m *Services) CheckUser(ctx context.Context, email string, contact string) error {
	ret := _m.Called(ctx, email, contact)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, contact)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckVenue provides a mock function with given fields: ctx, name, contact, email
func (_m *Services) CheckVenue(ctx context.Context, name string, contact string, email string) error {
	ret := _m.Called(ctx, name, contact, email)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, name, contact, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVenue provides a mock function with given fields: ctx, userID, id
func (_m *Services) DeleteVenue(ctx context.Context, userID int, id int) error {
	ret := _m.Called(ctx, userID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, userID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllBookings provides a mock function with given fields: ctx, userId
func (_m *Services) GetAllBookings(ctx context.Context, userId int) ([]*db.Booking, error) {
	ret := _m.Called(ctx, userId)

	var r0 []*db.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]*db.Booking, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []*db.Booking); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*db.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllVenues provides a mock function with given fields: ctx
func (_m *Services) GetAllVenues(ctx context.Context) ([]*db.Venue, error) {
	ret := _m.Called(ctx)

	var r0 []*db.Venue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*db.Venue, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*db.Venue); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*db.Venue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBooking provides a mock function with given fields: ctx, bookingid
func (_m *Services) GetBooking(ctx context.Context, bookingid int) (*db.Booking, error) {
	ret := _m.Called(ctx, bookingid)

	var r0 *db.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*db.Booking, error)); ok {
		return rf(ctx, bookingid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *db.Booking); ok {
		r0 = rf(ctx, bookingid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, bookingid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVenue provides a mock function with given fields: ctx, id
func (_m *Services) GetVenue(ctx context.Context, id int) (*db.Venue, error) {
	ret := _m.Called(ctx, id)

	var r0 *db.Venue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*db.Venue, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *db.Venue); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.Venue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoginUser provides a mock function with given fields: ctx, email, password
func (_m *Services) LoginUser(ctx context.Context, email string, password string) (string, error) {
	ret := _m.Called(ctx, email, password)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, email, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, email, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterUser provides a mock function with given fields: ctx, user
func (_m *Services) RegisterUser(ctx context.Context, user *db.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *db.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateVenue provides a mock function with given fields: ctx, venue, userID, id
func (_m *Services) UpdateVenue(ctx context.Context, venue *db.Venue, userID int, id int) error {
	ret := _m.Called(ctx, venue, userID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *db.Venue, int, int) error); ok {
		r0 = rf(ctx, venue, userID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewServices interface {
	mock.TestingT
	Cleanup(func())
}

// NewServices creates a new instance of Services. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewServices(t mockConstructorTestingTNewServices) *Services {
	mock := &Services{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
