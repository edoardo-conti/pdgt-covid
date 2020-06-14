package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type response struct {
	Status   int      `json:"status"`
	Messagge []string `json:"message"`
}

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
	//(todo) verificare se necessario: c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Benvenuto su PDGT-COVID!",
		"author":  "Edoardo Conti [278717]",
	})
}

// HandleAndamento ...
func HandleAndamento(c *gin.Context) {
	str := `{"status": 400, "message": ["/andamento/nazionale", "/andamento/regionale"]}`
	res := response{}
	json.Unmarshal([]byte(str), &res)

	c.JSON(http.StatusBadRequest, res)
}

//HandleNoRoute ...
func HandleNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "risorsa non disponibile",
	})
}
