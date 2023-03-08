package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/amancooks08/BookMySport/domain"
	"github.com/amancooks08/BookMySport/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mpatch "github.com/undefinedlabs/go-mpatch"
)

type CustomerHandlerTestSuite struct {
	suite.Suite
	service *mocks.Services
}

func TestCustomerHandlerTestSuite(t *testing.T) {
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func(patch *mpatch.Patch) {
		if err := patch.Unpatch(); err != nil {
			t.Fatal(err)
		}
	}(patch)

	suite.Run(t, new(CustomerHandlerTestSuite))
}

func (suite *CustomerHandlerTestSuite) SetupTest() {
	suite.service = &mocks.Services{}
}

func (suite *CustomerHandlerTestSuite) TearDownTest() {
	suite.service.AssertExpectations(suite.T())
}

func (suite *CustomerHandlerTestSuite) TestBookSlot() {
	t := suite.T()

	id := "id"
	t.Run("when valid booking request is made", func(t *testing.T) {
		reqBody := `{"booking_date":"2023-03-07","start_time":"21:00:00","end_time":"22:00:00","game":"Tennis"}`
		req := httptest.NewRequest(http.MethodPost, "/customer/venues/book?venueID=1", strings.NewReader(reqBody))
		rw := httptest.NewRecorder()
		ctx := req.Context()

		req = req.WithContext(context.WithValue(ctx, id, 1))

		respBody := domain.BookingResponse{
			Message: "Booking Successful.",
			Amount:  roundFloat(1000, 2),
		}

		reqBooking := domain.Booking{
			ID:          0,
			CustomerID:  1,
			VenueID:     1,
			BookingDate: "2023-03-07",
			BookingTime: time.Now().Format("2006-01-02 15:04:05.999999-07"),
			StartTime:   "21:00:00",
			EndTime:     "22:00:00",
			Game:        "Tennis",
			AmountPaid:  0,
		}

		suite.service.On("BookSlot", req.Context(), reqBooking).Return(1000.00, nil).Once()
		deps := dependencies{
			CustomerServices: suite.service,
		}

		exp, _ := json.Marshal(respBody)

		got := BookSlot(deps.CustomerServices)
		got.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, string(exp), rw.Body.String())

	})

	// t.Run("when date venueID is not provided", func(t *testing.T) {
	// 	reqBody := `{"booking_date":"2023-03-07","start_time":"21:00:00","end_time":"22:00:00","game":"Tennis"}`
	// 	req := httptest.NewRequest(http.MethodPost, "/customer/venues/book", strings.NewReader(reqBody))
	// 	rw := httptest.NewRecorder()
	// 	ctx := req.Context()

	// 	req = req.WithContext(context.WithValue(ctx, id, 1))

	// 	respBody := domain.Message{
	// 		Message: "please enter a venue ID",
	// 	}

	// 	reqBooking := domain.Booking{
	// 		ID:          0,
	// 		CustomerID:  1,
	// 		VenueID:     -1,
	// 		BookingDate: "2023-03-07",
	// 		BookingTime: time.Now().Format("2006-01-02 15:04:05.999999-07"),
	// 		StartTime:   "21:00:00",
	// 		EndTime:     "22:00:00",
	// 		Game:        "Tennis",
	// 		AmountPaid:  0,
	// 	}

	// 	suite.service.On("BookSlot", ctx, reqBooking).Return(1000.00, nil).Once()
	// 	deps := dependencies{
	// 		CustomerServices: suite.service,
	// 	}

	// 	exp, _ := json.Marshal(respBody)

	// 	got := BookSlot(deps.CustomerServices)
	// 	got.ServeHTTP(rw, req)

	// 	assert.Equal(t, http.StatusBadRequest, rw.Code)
	// 	assert.Equal(t, string(exp), rw.Body.String())

	// })

}
