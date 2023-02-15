package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		header := req.Header.Get("Authorization")

        // Check if the header is missing or invalid
        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            http.Error(rw, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Parse the JWT token from the header
        token, err := jwt.Parse(strings.TrimPrefix(header, "Bearer "), func(token *jwt.Token) (interface{}, error) {
            // Check the signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }

            // Set the secret key for the token
            return []byte("your-256-bit-secret"), nil
        })

        // Check if there was an error parsing the token
        if err != nil {
            http.Error(rw, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Check if the token is valid and has not expired
        if !token.Valid {
            http.Error(rw, "Unauthorized", http.StatusUnauthorized)
            return
        }
		// Get the claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            http.Error(rw, "Unauthorized", http.StatusUnauthorized)
            return
        }

		// Get the role from the claims
		role, ok := claims["role"].(string)
        if !ok {
            http.Error(rw, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Check if the user has the "admin" role to access admin routes
        if strings.HasPrefix(req.URL.Path, "/admin") {
            if role != "admin" {
                http.Error(rw, "Forbidden", http.StatusForbidden)
                return
            }
        }

		// Check if the user has the "customer" role to access customer routes
        if strings.HasPrefix(req.URL.Path, "/customer") {
            if role != "customer" {
                http.Error(rw, "Forbidden", http.StatusForbidden)
                return
            }
        }

		ctx := context.WithValue(req.Context(), "role", role)
		ctx = context.WithValue(ctx, "id", claims["user_id"])
        req = req.WithContext(ctx)
		
		// Call the next handler in the chain
		next.ServeHTTP(rw, req)
	})
}
