package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amancooks08/BookMySport/domain"
	"github.com/amancooks08/BookMySport/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	service *mocks.Services
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.service = &mocks.Services{}
}

func (suite *UserHandlerTestSuite) TearDownTest() {
	suite.service.AssertExpectations(suite.T())
}

// func (suite *UserHandlerTestSuite) TestPingHandler() {
// 	rw := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

// 	suite.service.On("Ping", mock.Anything).Return(nil)

// 	PingHandler(rw, req)

// 	suite.Equal(http.StatusOK, rw.Code)

// 	suite.service.AssertExpectations(suite.T())

// }

// func (suite *UserHandlerTestSuite) TestRegisterCustomer() {
// 	rw := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/register", nil)

// 	suite.service.On("RegisterCustomer", mock.Anything).Return(nil)

// 	RegisterCustomer(rw, req)

// 	suite.Equal(http.StatusOK, rw.Code)

// 	suite.service.AssertExpectations(suite.T())

// }

func (suite *UserHandlerTestSuite) TestLoginUser() {
	t := suite.T()
	t.Run("when valid user request is made", func(t *testing.T) {
		reqBody := `{"email":"cu1@gmail.com","password":"Password@123"}`
		req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
		rw := httptest.NewRecorder()
		ctx := req.Context()

		requestBody := domain.UserLogin{
			Email:    "cu1@gmail.com",
			Password: "Password@123",
		}

		responseBody := domain.LoginResponse{
			Token:   "token",
			Message: "Login Successful",
		}

		suite.service.On("LoginUser", ctx, requestBody.Email, requestBody.Password).Return("token", nil)
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := LoginUser(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})

	t.Run("when wrong password request is made", func(t *testing.T) {
		reqBody := `{"email":"cu1@gmail.com","password":"Password@321"}`
		req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
		rw := httptest.NewRecorder()
		ctx := req.Context()

		requestBody := domain.UserLogin{
			Email:    "cu1@gmail.com",
			Password: "Password@321",
		}

		responseBody := domain.Message{
			Message: "invalid credentials",
		}

		suite.service.On("LoginUser", ctx, requestBody.Email, requestBody.Password).Return("", errors.New("invalid credentials"))
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := LoginUser(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusUnauthorized, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})
}

func (suite *UserHandlerTestSuite) TestGetVenues() {
	t := suite.T()
	t.Run("when valid request is made to get all venues", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/venues", nil)
		rw := httptest.NewRecorder()
		ctx := req.Context()

		responseBody := []domain.Venue{
			{
				ID:          1,
				Name:        "Venue1",
				Address:     "Address1",
				City:        "City1",
				State:       "State1",
				Contact:     "1234567890",
				Email:       "email1@gmail.com",
				Opening:     "10:00",
				Closing:     "20:00",
				Price:       1000,
				Games: 	 	 []string{"Cricket", "Football"},
				OwnerID:     1,
			},

			{
				ID:          2,
				Name: 	  	 "Venue2",
				Address:     "Address2",
				City:        "City2",
				State:       "State2",
				Contact:     "1234567890",
				Email:       "email2@gmail.com",
				Opening:     "10:00",
				Closing:     "20:00",
				Price:       1000,
				Games: 	 	 []string{"Cricket", "Football"},
				OwnerID:     2,
			},
		}
		suite.service.On("GetAllVenues", ctx).Return([]domain.Venue{
			{
				ID:          1,
				Name:        "Venue1",
				Address:     "Address1",
				City:        "City1",
				State:       "State1",
				Contact:     "1234567890",
				Email:       "email1@gmail.com",
				Opening:     "10:00",
				Closing:     "20:00",
				Price:       1000,
				Games: 	 	 []string{"Cricket", "Football"},
				OwnerID:     1, 
			},

			{
				ID:          2,
				Name: 	  	 "Venue2",
				Address:     "Address2",
				City:        "City2",
				State:       "State2",
				Contact:     "1234567890",
				Email:       "email2@gmail.com",
				Opening:     "10:00",
				Closing:     "20:00",
				Price:       1000,
				Games: 	 	 []string{"Cricket", "Football"},
				OwnerID:     2,
			},

		}, nil).Once()
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := GetVenues(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})

	t.Run("when no venues are present and all venues were fetched", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/venues", nil)
		rw := httptest.NewRecorder()
		ctx := req.Context()

		responseBody := domain.Message{
			Message: "no venues found",
		}

		suite.service.On("GetAllVenues", ctx).Return([]domain.Venue{}, errors.New("no venues found"))
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := GetVenues(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})

	t.Run("when valid request is made to get a venue by its id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/venues?venueID=1", nil)
		rw := httptest.NewRecorder()
		ctx := req.Context()

		responseBody := domain.Venue{
			ID:          1,
			Name:        "Venue1",
			Address:     "Address1",
			City:        "City1",
			State:       "State1",
			Contact:     "1234567890",
			Email:       "email1@gmail.com",
			Opening:     "10:00",
			Closing:     "20:00",
			Price:       1000,
			Games: 	 	 []string{"Cricket", "Football"},
			OwnerID:     1,
		}

		suite.service.On("GetVenue", ctx, 1).Return(domain.Venue{
			ID:          1,
			Name:        "Venue1",
			Address:     "Address1",
			City:        "City1",
			State:       "State1",
			Contact:     "1234567890",
			Email:       "email1@gmail.com",
			Opening:     "10:00",
			Closing:     "20:00",
			Price:       1000,
			Games: 	 	 []string{"Cricket", "Football"},
			OwnerID:     1,
		}, nil).Once()

		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := GetVenues(deps.CustomerServices)
		got.ServeHTTP(rw, req)
		
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})

	t.Run("when invalid request is made to get a venue by its id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/venues?venueID=1", nil)
		rw := httptest.NewRecorder()
		ctx := req.Context()

		responseBody := domain.Message{
			Message: "no venue found",
		}

		suite.service.On("GetVenue", ctx, 1).Return(domain.Venue{}, errors.New("no venue found"))
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := GetVenues(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})
}
func TestPingHandler(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PingHandler(tt.args.rw, tt.args.req)
		})
	}
}


func (suite *UserHandlerTestSuite) TestCheckAvailabililty(){
	t := suite.T()

	t.Run("when valid request is made to check availability of a venue", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/venues/slots?venueID=1&date=2023-03-31", nil)
		rw := httptest.NewRecorder()
		ctx := req.Context()

		venueID := GetVenueID(req)
		date := req.URL.Query().Get("date")
		responseBody := []domain.Slot{
			{
				VenueID: 1,
				StartTime: "10:00",
				EndTime: "11:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "11:00",
				EndTime: "12:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "12:00",
				EndTime: "13:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "13:00",
				EndTime: "14:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "14:00",
				EndTime: "15:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "15:00",
				EndTime: "16:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "16:00",
				EndTime: "17:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "17:00",
				EndTime: "18:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "18:00",
				EndTime: "19:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "19:00",
				EndTime: "20:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "20:00",
				EndTime: "21:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "21:00",
				EndTime: "22:00",
				Date: "2023-03-31",
			},
		}

		suite.service.On("CheckAvailability", ctx, venueID, date).Return([]domain.Slot{
			{
				VenueID: 1,
				StartTime: "10:00",
				EndTime: "11:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "11:00",
				EndTime: "12:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "12:00",
				EndTime: "13:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "13:00",
				EndTime: "14:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "14:00",
				EndTime: "15:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "15:00",
				EndTime: "16:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "16:00",
				EndTime: "17:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "17:00",
				EndTime: "18:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "18:00",
				EndTime: "19:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "19:00",
				EndTime: "20:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "20:00",
				EndTime: "21:00",
				Date: "2023-03-31",
			},

			{
				VenueID: 1,
				StartTime: "21:00",
				EndTime: "22:00",
				Date: "2023-03-31",
			},
		}, nil).Once()
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := CheckAvailability(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})
}