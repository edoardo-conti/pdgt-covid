package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/edoardo-conti/pdgt-covid/models"
	
	// importare: "github.com/edoardo-conti/pdgt-covid/models"
)

func nationalTrend(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM nazione")
	if err != nil {
		log.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	//got := []nazione{}

	var nazioni []models.nationalTrend
	for rows.Next() {
		// c := new(Course)
		var r models.nationalTrend
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
		nazioni = append(nazioni, nationalTrend{r.Data, r.Stato, r.RicoveratiConSintomi, r.TerapiaIntensiva, r.TotaleOspedalizzati, r.IsolamentoDomiciliare, r.TotalePositivi, r.VariazioneTotalePositivi, r.NuoviPositivi, r.DimessiGuariti, r.Deceduti, r.TotaleCasi, r.Tamponi, r.CasiTestati, r.NoteIT, r.NoteEN})

		//got = append(got, r)
	}

	//log.Println(got)
	//nazioniBytes, _ := json.Marshal(&nazioni)

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": nazioni,
	})
}

func nationalTrendByDate(c *gin.Context) {
	nazioneDate := c.Params.ByName("bydate")

	if nazioneDate != "" {
		data := "2020-02-28"
		var r nationalTrend

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
}
