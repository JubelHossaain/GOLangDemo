package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the PostgreSQL database connection
func InitDB() {
	var err error
	//connStr := "postgres://postgres:admin@localhost:6061/postgres?sslmode=disable"

	connStr := "postgresql://jubel:SpmUcOUGSzWOk0zlVyQK2kV8yFPi2oXA@dpg-d02ccp3e5dus73bku370-a/demo_user_db"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	 // Create users table
	 _, err = DB.Exec(`
	 CREATE TABLE IF NOT EXISTS users (
		 id SERIAL PRIMARY KEY,
		 username TEXT NOT NULL UNIQUE,
		 password TEXT NOT NULL
	 );
 `)
 if err != nil {
	 log.Fatal("Error creating users table:", err)
 }

 // Create messages table with foreign key
 _, err = DB.Exec(`
	 CREATE TABLE IF NOT EXISTS messages (
		 id SERIAL PRIMARY KEY,
		 user_id INTEGER NOT NULL,
		 content TEXT NOT NULL,
		 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	 );
 `)
 if err != nil {
	 log.Fatal("Error creating messages table:", err)
 }

	fmt.Println("Connected to the database successfully!")
}
