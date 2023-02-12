package service

import (
	"encoding/json"
	"net/http"
	"regexp"

	db "github.com/amancooks08/BookMySport/db"
	"golang.org/x/crypto/bcrypt"
)

func registerUser(rw http.ResponseWriter, req *http.Request, deps dependencies, userType string) {
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
	if user.Name == "" || user.Contact == "" || user.Email == "" || user.City == "" || user.State == "" {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if validateContact(user.Contact) && validateEmail(user.Email) {

		// Store the user in the database
		// var cu *db.Customer
		user.Type = userType
		err = deps.CustomerServices.RegisterUser(req.Context(), &user)
		if err != nil {
			http.Error(rw, "Failed to register user", http.StatusInternalServerError)
			return
		}
		registerUserResponse := &User{
			Id:      user.Id,
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
			http.Error(rw, "Failed to marshal response", http.StatusInternalServerError)
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
