package models

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

func ConnectDatabase() {
	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// bind variabili
	DB = database
}
