package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// repository contains the details of a repository
type entrySummary struct {
	data                     string
	stato                    string
	ricoveratiConSintomi     int
	terapiaIntensiva         int
	totaleOspedalizzati      int
	isolamentoDomiciliare    int
	totalePositivi           int
	variazioneTotalePositivi int
	nuoviPositivi            int
	dimessiGuariti           int
	deceduti                 int
	totaleCasi               int
	tamponi                  int
	casiTestati              int
	noteIT                   string
	noteEN                   string
}

type entries struct {
	Entries []entrySummary
}

// queryEntries first fetches the repositories data from the db
func queryEntries(righe *entries, db *sql.DB) error {
	rows, err := db.Query(`SELECT * FROM nazione`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		riga := entrySummary{}
		err = rows.Scan(
			&riga.data,
			&riga.stato,
			&riga.ricoveratiConSintomi,
			&riga.terapiaIntensiva,
			&riga.totaleOspedalizzati,
			&riga.isolamentoDomiciliare,
			&riga.totalePositivi,
			&riga.variazioneTotalePositivi,
			&riga.nuoviPositivi,
			&riga.dimessiGuariti,
			&riga.deceduti,
			&riga.totaleCasi,
			&riga.tamponi,
			&riga.casiTestati,
			&riga.noteIT,
			&riga.noteEN,
		)
		if err != nil {
			return err
		}
		righe.Entries = append(righe.Entries, riga)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// db
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	// endpoint: /
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"covid-19 api": "welcome"})
	})

	// endpoint: /nazione
	router.GET("/nazione", func(c *gin.Context) {
		righe := entries{}

		err := queryEntries(&righe, db)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error 1: ", err))
			return
		}

		out, err := json.Marshal(righe)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error 2: ", err))
			return
		}

		c.JSON(http.StatusOK, string(out))
	})

	router.Run(":" + port)
}
