package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amancooks08/BookMySport/db"
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

		requestBody := UserLogin{
			Email:    "cu1@gmail.com",
			Password: "Password@123",
		}

		responseBody := LoginResponse{
			Token:   "token",
			Message: "Login Successful",
		}

		suite.service.On("LoginUser", ctx, requestBody).Return("token", nil)
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

		requestBody := UserLogin{
			Email:    "cu1@gmail.com",
			Password: "Password@321",
		}

		responseBody := Message{
			Message: "invalid credentials",
		}

		suite.service.On("LoginUser", ctx, requestBody).Return("", errors.New("invalid credentials"))
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

func (suite *UserHandlerTestSuite) TestGetAllVenues() {
	t := suite.T()
	t.Run("when valid request is made", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/venues", nil)
		rw := httptest.NewRecorder()
		ctx := req.Context()

		responseBody := []*db.Venue{
			{
				ID:      1,
				Name:    "Venue1",
				Address: "Address1",
				City:    "City1",
				State:   "State1",
				Contact: "1234567890",
				Email:   "email1@gmail.com",
				Opening: "06:00",
				Closing: "23:00",
				Price:   1000,
				Games:   []string{"Cricket", "Football"},
				Rating:  4.5,
			},

			{
				ID:      2,
				Name:    "Venue2",
				Address: "Address2",
				City:    "City2",
				State:   "State2",
				Contact: "9383747727",
				Email:   "email2@gmail.com",
				Opening: "06:00",
				Closing: "23:00",
				Price:   1100,
				Games:   []string{"Cricket", "Football"},
				Rating:  4.5,
			},
		}
		suite.service.On("GetAllVenues", ctx).Return([]*db.Venue{
			{
				ID:      1,
				Name:    "Venue1",
				Address: "Address1",
				City:    "City1",
				State:   "State1",
				Contact: "1234567890",
				Email:   "email1@gmail.com",
				Opening: "06:00",
				Closing: "23:00",
				Price:   1000,
				Games:   []string{"Cricket", "Football"},
				Rating:  4.5,
			},
			{
				ID:      2,
				Name:    "Venue2",
				Address: "Address2",
				City:    "City2",
				State:   "State2",
				Contact: "9383747727",
				Email:   "email2@gmail.com",
				Opening: "06:00",
				Closing: "23:00",
				Price:   1100,
				Games:   []string{"Cricket", "Football"},
				Rating:  4.5,
			},
		}, nil)
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := GetAllVenues(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())
	})

	t.Run("when no venues are present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/venues", nil)
		rw := httptest.NewRecorder()
		ctx := req.Context()

		responseBody := Message{
			Message: "No Venues Found",
		}

		suite.service.On("GetAllVenues", ctx).Return([]*db.Venue{}, errors.New("No Venues Found"))
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(responseBody)

		got := GetAllVenues(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
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
