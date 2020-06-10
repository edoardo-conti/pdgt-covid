package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// repository contains the details of a repository
type nazione struct {
	Data                     string         `json:"data"`
	Stato                    string         `json:"stato"`
	RicoveratiConSintomi     int            `json:"ricoverati_con_sintomi"`
	TerapiaIntensiva         int            `json:"terapia_intensiva"`
	TotaleOspedalizzati      int            `json:"totale_ospedalizzati"`
	IsolamentoDomiciliare    int            `json:"isolamento_domiciliare"`
	TotalePositivi           int            `json:"totale_oositivi"`
	VariazioneTotalePositivi int            `json:"variazione_totale_positivi"`
	NuoviPositivi            int            `json:"nuovi_positivi"`
	DimessiGuariti           int            `json:"dimessi_guariti"`
	Deceduti                 int            `json:"deceduti"`
	TotaleCasi               int            `json:"totale_casi"`
	Tamponi                  int            `json:"tamponi"`
	CasiTestati              sql.NullInt64  `json:"casi_testati"`
	NoteIT                   sql.NullString `json:"note_it"`
	NoteEN                   sql.NullString `json:"note_en"`
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// db
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())

	// endpoint: /
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Benvenuto su PDGT-COVID!",
		})
	})

	// endpoint: /nazione
	router.GET("/nazione", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM nazione")
		if err != nil {
			log.Fatalf("Query: %v", err)
		}
		defer rows.Close()

		//got := []nazione{}

		var nazioni []nazione
		for rows.Next() {
			// c := new(Course)
			var r nazione
			err = rows.Scan(
				&r.Data,
				&r.Stato,
				&r.RicoveratiConSintomi,
				&r.TerapiaIntensiva,
				&r.TotaleOspedalizzati,
				&r.IsolamentoDomiciliare,
				&r.TotalePositivi,
				&r.VariazioneTotalePositivi,
				&r.NuoviPositivi,
				&r.DimessiGuariti,
				&r.Deceduti,
				&r.TotaleCasi,
				&r.Tamponi,
				&r.CasiTestati,
				&r.NoteIT,
				&r.NoteEN)
			if err != nil {
				log.Fatalf("Scan: %v", err)
			}
			nazioni = append(nazioni, nazione{r.Data, r.Stato, r.RicoveratiConSintomi, r.TerapiaIntensiva, r.TotaleOspedalizzati, r.IsolamentoDomiciliare, r.TotalePositivi, r.VariazioneTotalePositivi, r.NuoviPositivi, r.DimessiGuariti, r.Deceduti, r.TotaleCasi, r.Tamponi, r.CasiTestati, r.NoteIT, r.NoteEN})

			//got = append(got, r)
		}

		//log.Println(got)
		//nazioniBytes, _ := json.Marshal(&nazioni)

		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": nazioni,
		})
	})

	// endpoint: /nazione/:bydate (not work)
	router.GET("/nazione/:bydate", func(c *gin.Context) {
		nazioneDate := c.Params.ByName("bydate")

		if nazioneDate != "" {
			data := "2020-02-28"
			var r nazione

			row := db.QueryRow("SELECT * FROM nazione WHERE data=$1", data)
			switch err := row.Scan(
				&r.Data,
				&r.Stato,
				&r.RicoveratiConSintomi,
				&r.TerapiaIntensiva,
				&r.TotaleOspedalizzati,
				&r.IsolamentoDomiciliare,
				&r.TotalePositivi,
				&r.VariazioneTotalePositivi,
				&r.NuoviPositivi,
				&r.DimessiGuariti,
				&r.Deceduti,
				&r.TotaleCasi,
				&r.Tamponi,
				&r.CasiTestati,
				&r.NoteIT,
				&r.NoteEN); err {
			case sql.ErrNoRows:
				fmt.Println("No rows were returned!")
			case nil:

				c.JSON(http.StatusOK, gin.H{
					"status":  200,
					"message": r,
				})

			default:
				panic(err)
			}
		}

	})

	router.Run(":" + port)
}
