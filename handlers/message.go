package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"GOFolder/models"
	"GOFolder/storage"
)

// CreateMessageHandler handles message creation
func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var msg models.Message

		// Parse the JSON body
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Insert the message into the database
		query := "INSERT INTO messages (user_id, content) VALUES ($1, $2) RETURNING id"
		err = storage.DB.QueryRow(query, msg.UserID, msg.Content).Scan(&msg.ID)
		if err != nil {
			http.Error(w, "Error creating message: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the created message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ReadMessageHandler handles fetching messages for a user
func ReadMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// Query the database for messages
		query := "SELECT id, user_id, content FROM messages WHERE user_id = $1"
		rows, err := storage.DB.Query(query, userID)
		if err != nil {
			http.Error(w, "Error fetching messages: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse the results
		var messages []models.Message
		for rows.Next() {
			var msg models.Message
			err := rows.Scan(&msg.ID, &msg.UserID, &msg.Content)
			if err != nil {
				http.Error(w, "Error parsing messages: "+err.Error(), http.StatusInternalServerError)
				return
			}
			messages = append(messages, msg)
		}

		// Respond with the messages
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UpdateMessageHandler handles updating an existing message
func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg models.Message

	// Parse the JSON body
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate input
	if msg.ID == 0 || msg.UserID == 0 || msg.Content == "" {
		http.Error(w, "Message ID, User ID, and Content are required", http.StatusBadRequest)
		return
	}

	// Update the message in the database
	query := "UPDATE messages SET content = $1 WHERE id = $2 AND user_id = $3"
	result, err := storage.DB.Exec(query, msg.Content, msg.ID, msg.UserID)
	if err != nil {
		http.Error(w, "Error updating message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "No message found to update", http.StatusNotFound)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message updated successfully"))
}

// DeleteMessageHandler handles deleting a message
func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the message ID and user ID from the query parameters
	messageIDStr := r.URL.Query().Get("id")
	userIDStr := r.URL.Query().Get("user_id")

	if messageIDStr == "" || userIDStr == "" {
		http.Error(w, "Message ID and User ID are required", http.StatusBadRequest)
		return
	}

	// Convert the IDs to integers
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, "Invalid Message ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Delete the message from the database
	query := "DELETE FROM messages WHERE id = $1 AND user_id = $2"
	result, err := storage.DB.Exec(query, messageID, userID)
	if err != nil {
		http.Error(w, "Error deleting message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "No message found to delete", http.StatusNotFound)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message deleted successfully"))
}
