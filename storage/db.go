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

	fmt.Println("Connected to the database successfully!")
}
