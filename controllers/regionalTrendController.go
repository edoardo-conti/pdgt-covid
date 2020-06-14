package controllers

import (
	"log"
	"net/http"
	"pdgt-covid/models"

	"github.com/gin-gonic/gin"
)

var queryFirstPart string = `SELECT data, json_agg(json_build_object('stato', stato,
							'codice_regione', codice_regione,
							'denominazione_regione', denominazione_regione,
							'lat', lat,
							'long', long,
							'ricoverati_con_sintomi', ricoverati_con_sintomi,
							'terapia_intensiva', terapia_intensiva,
							'totale_ospedalizzati', totale_ospedalizzati,
							'isolamento_domiciliare', isolamento_domiciliare,
							'totale_positivi', totale_positivi,
							'variazione_totale_positivi', variazione_totale_positivi,
							'nuovi_positivi', nuovi_positivi,
							'dimessi_guariti', dimessi_guariti,
							'deceduti', deceduti,
							'totale_casi', totale_casi,
							'tamponi', tamponi,
							'casi_testati', casi_testati,
							'note_it', note_it,
							'note_en', note_en)) as records 
							FROM regioni`

// RegionalTrend ...
func RegionalTrend(c *gin.Context) {
	// todo: idea per multiplexare richieste differenti con una singola query a stringhe composte
	query := queryFirstPart + " GROUP BY regioni.data ORDER BY regioni.data ASC;"
	rows, err := models.DB.Query(query)
	if err != nil {
		log.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	// counter per il conteggio dei record
	counter := 0

	var regioniCollect []models.RegionalTrendCollect
	for rows.Next() {
		var r models.RegionalTrendCollect

		err = rows.Scan(&r.Data, &r.Info)
		if err != nil {
			log.Fatalf("scan error: %v", err)
		}

		regioniCollect = append(regioniCollect, models.RegionalTrendCollect{r.Data, r.Info})

		counter++
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"count":  counter,
		"data":   regioniCollect,
	})
}
