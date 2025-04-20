package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for signing JWT tokens
// var jwtKey = []byte("your_secret_key") // Replace with a secure and private key
var jwtKey = []byte("wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY")

// ValidateToken middleware validates the JWT token and extracts the user ID
func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Check if the token is blacklisted
		tokenBlacklist.RLock()
		if tokenBlacklist.tokens[tokenString] {
			tokenBlacklist.RUnlock()
			http.Error(w, "Token is invalid or expired", http.StatusUnauthorized)
			return
		}
		tokenBlacklist.RUnlock()

		// Parse and validate the token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract the user ID from the claims
		userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
		if !ok {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request header as a string
		r.Header.Set("user_id", strconv.Itoa(int(userID)))

		// Call the next handler
		next(w, r)
	}
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID int) (string, error) {
	// Define the claims for the token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtKey)
}
