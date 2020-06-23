package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// response struttura impiegata per formulare i messaggi di risposta JSON
type response struct {
	Status   int      `json:"status"`
	Messagge []string `json:"message"`
}

// DB variabile "globale" per la condivisione della connessione al database
var DB *sql.DB

// ConnectDatabase funzione per connttere il progetto al database Heroku Postegres
func ConnectDatabase() {
	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// bind variabili
	DB = database
}

// HandleWelcome handler dell'endpoint root '/'
func HandleWelcome(c *gin.Context) {
	/*
	 * gin si occupa automaticamente di specificare Content-Type e encoding della risposta:
	 * c.Header("Content-Type", "application/json")
	 */
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Benvenuto su pdgt-covid!",
		"author":  "Edoardo Conti [278717]",
	})
}

// HandleAndamento handler dell'endpoint informativo '/andamento'
func HandleAndamento(c *gin.Context) {
	str := `{"status": 400, "message": ["/trend/nazionale", "/trend/regionale"]}`
	res := response{}

	// unmarshal del messaggio informativo
	json.Unmarshal([]byte(str), &res)

	c.JSON(http.StatusBadRequest, res)
}

//HandleNoRoute handler per evitare inconsistenza del web service riguardo a richieste senza endpoint
func HandleNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Risorsa non disponibile.",
	})
}
