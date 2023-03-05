package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/amancooks08/BookMySport/domain"
	"github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func registerUser(rw http.ResponseWriter, req *http.Request, CustomerServices Services, userType string) {
	if req.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
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
		msg := domain.Message{
			Message: "please don't leave any field empty",
		}
		json_response, _ := json.Marshal(msg)

		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(json_response)
		return
	}

	if validateContact(user.Contact) && validateEmail(user.Email) {

		user.Type = userType
		err = CustomerServices.RegisterUser(req.Context(), user)
		if err != nil {
			http.Error(rw, "Failed to register user", http.StatusInternalServerError)
			return
		}
		registerUserResponse := domain.User{
			ID:      user.ID,
			Name:    user.Name,
			Contact: user.Contact,
			Email:   user.Email,
			City:    user.City,
			State:   user.State,
			Type:    userType,
		}
		// Send a successful response

		json_response, err := json.Marshal(registerUserResponse)
		if err != nil {
			msg := domain.Message{
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
		msg := domain.Message{
			Message: "invalid contact or email details.",
		}
		json_response, _ := json.Marshal(msg)

		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(json_response)
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

func GetUserID(req *http.Request, rw http.ResponseWriter) (int) {
	header := req.Header.Get("Authorization")

	// Check if the header is missing or invalid
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		http.Error(rw, "Unauthorized7", http.StatusUnauthorized)
		return 0
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
		return 0
	}

	// Check if the token is valid and has not expired
	if !token.Valid {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return 0
	}
	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		http.Error(rw, "Token Invalid", http.StatusUnauthorized)
		return 0
	}

	userID, ok := claims["user_id"].(float64)
	if !ok || !token.Valid {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return 0
	}
	return int(userID)
}

func GetVenueID(req *http.Request) int {
	if(req.URL.Query().Get("venueID") == ""){
		return 0
	}
	venueID, err := strconv.Atoi(req.URL.Query().Get("venueID"))
	if err != nil {
		logger.WithField("error", err).Error("error while parsing venueID")
		return 0
	}

	return venueID
}

type dependencies struct {
	CustomerServices Services
}
