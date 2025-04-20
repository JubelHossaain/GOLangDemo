package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"GOFolder/models"
	"GOFolder/storage"
)

var mutex = &sync.Mutex{}

// SignupHandler handles user signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User

	// Parse the JSON body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate the input
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err = storage.DB.QueryRow(query, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials models.User
	var user models.User

	// Parse the JSON body
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Query the database for the user
	query := "SELECT id, password FROM users WHERE username = $1"
	err = storage.DB.QueryRow(query, credentials.Username).Scan(&user.ID, &user.Password)
	if err != nil {
		log.Printf("Error querying user: %v", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check if the password matches
	if user.Password != credentials.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a token for the user
	token, err := GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Respond with the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

var tokenBlacklist = struct {
	sync.RWMutex
	tokens map[string]bool
}{tokens: make(map[string]bool)}

// SignOutHandler handles user sign-out
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
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

	// Add the token to the blacklist
	tokenBlacklist.Lock()
	tokenBlacklist.tokens[tokenString] = true
	tokenBlacklist.Unlock()

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signed out successfully"))
}
