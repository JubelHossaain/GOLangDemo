package main

import (
	"fmt"
	"net/http"

	"GOFolder/handlers"
	"GOFolder/storage"
)

func main() {
	// Initialize the database
	storage.InitDB()

	// Serve templates
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/users", handlers.UsersPageHandler)
	http.HandleFunc("/messages", handlers.MessagesPageHandler)

	// API routes
	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/signout", handlers.ValidateToken(handlers.SignOutHandler))
	http.HandleFunc("/create", handlers.ValidateToken(handlers.CreateMessageHandler))
	http.HandleFunc("/read", handlers.ValidateToken(handlers.ReadMessageHandler))
	http.HandleFunc("/update", handlers.ValidateToken(handlers.UpdateMessageHandler))
	http.HandleFunc("/delete", handlers.ValidateToken(handlers.DeleteMessageHandler))

	// Start the server
	fmt.Println("Server is running on http://localhost:6060")
	http.ListenAndServe(":6060", nil)
}
