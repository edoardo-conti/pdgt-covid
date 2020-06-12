package models

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//DB ...
var DB *sql.DB

//ConnectDatabase ...
func ConnectDatabase() {
	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// bind variabili
	DB = database
}

//HandleWelcome ...
func HandleWelcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Benvenuto su PDGT-COVID!",
		"author":  "Edoardo Conti [278717]",
	})
}
