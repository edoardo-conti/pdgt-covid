package controllers

import (
	"log"
	"net/http"
	"pdgt-covid/models"

	"github.com/gin-gonic/gin"
)

// RegionalTrend ...
func RegionalTrend(c *gin.Context) {
	rows, err := models.DB.Query("SELECT * FROM regioni")
	if err != nil {
		log.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	// counter per il conteggio dei record
	counter := 0

	var regioni []models.RegionalTrend
	for rows.Next() {
		var r models.RegionalTrend
		err = rows.Scan(
			&r.Data,
			&r.Stato,
			&r.CodiceRegione,
			&r.DenominazioneRegione,
			&r.Lat,
			&r.Long,
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
			log.Fatalf("scan error: %v", err)
		}

		regioni = append(regioni, models.RegionalTrend{r.Data, r.Stato, r.CodiceRegione, r.DenominazioneRegione, r.Lat, r.Long, r.RicoveratiConSintomi, r.TerapiaIntensiva, r.TotaleOspedalizzati, r.IsolamentoDomiciliare, r.TotalePositivi, r.VariazioneTotalePositivi, r.NuoviPositivi, r.DimessiGuariti, r.Deceduti, r.TotaleCasi, r.Tamponi, r.CasiTestati, r.NoteIT, r.NoteEN})

		// incremento del counter
		counter++
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"count":  counter,
		"data":   regioni,
	})
}
