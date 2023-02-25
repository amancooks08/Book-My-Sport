package service

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	db "github.com/amancooks08/BookMySport/db"
// 	"github.com/stretchr/testify/assert"
// )

// func TestBookSlotHandler(t *testing.T) {
// 	// Test if the handler returns a "Method Not Allowed" error for a GET request
// 	req, err := http.NewRequest("GET", "/customer/venues/{venue_id}/slots", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(BookSlot(dependencies{}))
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusMethodNotAllowed {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusMethodNotAllowed)
// 	}

// 	// Test if the handler returns a "Method Not Allowed" error for a PUT request
// 	req, err = http.NewRequest("PUT", "/customer/venues/{venue_id}/slots", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr = httptest.NewRecorder()
// 	handler = http.HandlerFunc(BookSlot(dependencies{}))
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusMethodNotAllowed {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusMethodNotAllowed)
// 	}

// 	// Test if the handler parses the json Body correctly
// 	booking := db.Booking{
// 		BookedBy:    1,
// 		BookedAt:    1,
// 		BookingDate: time.Now().Format("2006-01-02"),
// 		BookingTime: time.Now().Format("2006-01-02 15:04:05.999999-07"),
// 		StartTime:   "10:00:00",
// 		EndTime:     "11:00:00",
// 		Game:        "Football",
// 		AmountPaid:  1000,
// 	}
// 	bookingJson, err := json.Marshal(booking)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err = http.NewRequest("POST", "/customer/venues/{venue_id}/slots", strings.NewReader(string(bookingJson)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr = httptest.NewRecorder()
// 	handler = http.HandlerFunc(BookSlot(dependencies{}))
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusCreated)
// 	}

// 	// Test if the handler returns a "Bad Request" error for an invalid request body
// 	req, err = http.NewRequest("POST", "/customer/venues/{venue_id}/slots", strings.NewReader("invalid json"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr = httptest.NewRecorder()
// 	handler = http.HandlerFunc(BookSlot(dependencies{}))
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusBadRequest)
// 	}

// 	// Test if the handler returns a "Bad Request" error for invalid details
// 	booking = db.Booking{
// 		BookedBy:    0,
// 		BookedAt:    0,
// 		BookingDate: "invalid date",
// 		BookingTime: "invalid time",
// 		StartTime:   "invalid start time",
// 		EndTime:     "invalid end time",
// 		Game:        "",
// 		AmountPaid:  0,
// 	}
// 	bookingJson, err = json.Marshal(booking)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err = http.NewRequest("POST", "/customer/venues/{venue_id}/slots", strings.NewReader(string(bookingJson)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr = httptest.NewRecorder()
// 	handler = http.HandlerFunc(BookSlot(dependencies{}))
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusBadRequest)
// 	}

// 	// Test if the handler returns a "Internal Server Error" error if the service layer returns an error
// 	booking = db.Booking{
// 		BookedBy:    1,
// 		BookedAt:    1,
// 		BookingDate: time.Now().Format("2006-01-02"),
// 		BookingTime: time.Now().Format("2006-01-02 15:04:05.999999-07"),
// 		StartTime:   "10:00:00",
// 		EndTime:     "11:00:00",
// 		Game:        "Football",
// 		AmountPaid:  1000,
// 	}
// 	tests := []struct {
// 		name           string
// 		mockFunc       func(ctx context.Context, booking *Booking) error // mock function to replace the actual CustomerServices.BookSlot function
// 		ctx            context.Context
// 		expectedStatus int
// 	}{
// 		{
// 			name: "Success",
// 			mockFunc: func(ctx context.Context, booking *Booking) error {
// 				return nil // mock a successful booking
// 			},
// 			ctx:            context.Background(),
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name: "Internal Server Error",
// 			mockFunc: func(ctx context.Context, booking *Booking) error {
// 				return errors.New("failed to book slot") // mock a failed booking
// 			},
// 			ctx:            context.Background(),
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// create a mock CustomerServices.BookSlot function that returns the error specified in mockFunc

// 			// create a mock request and response
// 			r, err := http.NewRequest(http.MethodPost, "/book", nil)
// 			assert.NoError(t, err)

// 			w := httptest.NewRecorder()

// 			// call the handler with the mock request and response
// 			BookSlot(w, r.WithContext(tt.ctx), booking, CustomerServices)

// 			// assert that the response status code matches the expected value
// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 		})
// 	}
// }
