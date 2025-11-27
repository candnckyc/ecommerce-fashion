package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to PostgreSQL
func Connect(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✓ Database connected successfully")
	return db
}
