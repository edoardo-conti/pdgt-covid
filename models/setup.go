package models

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func connectDatabase() {
	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// bind variabili
	db = database
}
