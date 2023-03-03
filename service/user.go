package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	db "github.com/amancooks08/BookMySport/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func registerUser(rw http.ResponseWriter, req *http.Request, CustomerServices Services, userType string) {
	if req.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user db.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err = CustomerServices.CheckUser(req.Context(), user.Email, user.Contact)
	if err != nil {
		http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if user.Name == "" || user.Contact == "" || user.Email == "" || user.City == "" || user.State == "" || user.Password == "" {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if validateContact(user.Contact) && validateEmail(user.Email) {

		// Store the user in the database
		// var cu *db.Customer
		user.Type = userType
		err = CustomerServices.RegisterUser(req.Context(), &user)
		if err != nil {
			http.Error(rw, "Failed to register user", http.StatusInternalServerError)
			return
		}
		registerUserResponse := &User{
			Id:      user.ID,
			Name:    user.Name,
			Contact: user.Contact,
			Email:   user.Email,
			City:    user.City,
			State:   user.State,
			Type:    userType,
		}
		// Send a successful response
		// rw.Write([]byte("User registered successfully"))
		json_response, err := json.Marshal(registerUserResponse)
		if err != nil {
			msg := Message{
				Message: "Failed to register user",
			}
			json_response, _ := json.Marshal(msg)

			rw.WriteHeader(http.StatusInternalServerError)
			rw.Header().Add("Content-Type", "application/json")
			rw.Write(json_response)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(json_response)
	} else {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
}

func validateEmail(email string) bool {
	// Define a shorter regular expression pattern for email addresses
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func validateContact(contact string) bool {
	// Define a shorter regular expression pattern for contact numbers
	re := regexp.MustCompile(`^[0-9]{10,}$`)
	return re.MatchString(contact)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GetUserVenueId(req *http.Request, rw http.ResponseWriter) (int, int) {
	header := req.Header.Get("Authorization")

	// Check if the header is missing or invalid
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return 0, 0
	}

	// Parse the JWT token from the header
	token, err := jwt.Parse(strings.TrimPrefix(header, "Bearer "), func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Set the secret key for the token
		return []byte("secret@987"), nil
	})

	// Check if there was an error parsing the token
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return 0, 0
	}

	// Check if the token is valid and has not expired
	if !token.Valid {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return 0, 0
	}
	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		http.Error(rw, "Token Invalid", http.StatusUnauthorized)
		return 0, 0
	}

	user_id, ok := claims["user_id"].(float64)
	if !ok || !token.Valid {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return 0, 0
	}
	vars := mux.Vars(req)
	venueID, err := strconv.Atoi(vars["venue_id"])
	if err != nil {
		logger.WithField("error", err).Error("Error while parsing venue_id")
		return int(user_id), 0
	}
	return int(user_id), venueID
}
