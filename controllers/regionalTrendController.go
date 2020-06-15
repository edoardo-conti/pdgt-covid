package controllers

import (
	"log"
	"net/http"
	"pdgt-covid/models"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
func IsValidRegion(region string) bool {
	switch region {
	case
		"Abruzzo",
		"Basilicata",
		"P.A. Bolzano",
		"Calabria",
		"Campania",
		"Emilia-Romagna",
		"Friuli Venezia Giulia",
		"Lazio",
		"Liguria",
		"Lombardia",
		"Marche",
		"Molise",
		"Piemonte",
		"Puglia",
		"Sardegna",
		"Sicilia",
		"Toscana",
		"P.A. Trento",
		"Umbria",
		"Veneto":
		return true
	}
	return false
}
*/

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

//RegionalTrendHandler ...
func RegionalTrendHandler(mode int) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// todo: idea per multiplexare richieste differenti con una singola query a stringhe composte
		query := ""
		if mode == 1 {
			// query NON filtrata
			query = queryFirstPart + " GROUP BY regioni.data ORDER BY regioni.data ASC;"
		} else if mode == 2 {
			// filtrato per data
			data := c.Params.ByName("bydata")
			// check validità parametro
			dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
			if data != "" && dateCheck.MatchString(data) {
				// valido
				query = queryFirstPart + " WHERE data='" + data + "' GROUP BY regioni.data ORDER BY regioni.data ASC;"
			}
		} else if mode == 3 {
			// filtrato per id regione
			regID := c.Params.ByName("byregid")
			// check validità parametro
			if regIDint, err := strconv.Atoi(regID); err == nil {
				if regIDint >= 0 && regIDint <= 22 {
					query = queryFirstPart + " WHERE codice_regione='" + regID + "' GROUP BY regioni.data ORDER BY regioni.data ASC;"
				}
			}
		} else if mode == 4 {
			// filtrato per picco
			query = queryFirstPart + " WHERE nuovi_positivi=(SELECT MAX(nuovi_positivi) FROM regioni) GROUP BY regioni.data ORDER BY regioni.data ASC;"
		} else if mode == 5 {
			// filtrato per picco di regione
			regID := c.Params.ByName("byregid")
			// check validità parametro
			if regIDint, err := strconv.Atoi(regID); err == nil {
				if regIDint >= 0 && regIDint <= 22 {
					query = queryFirstPart + " WHERE codice_regione='" + regID + "' AND nuovi_positivi=(select max(nuovi_positivi) from regioni where codice_regione='" + regID + "' group by codice_regione) GROUP BY regioni.data ORDER BY regioni.data ASC;"
				}
			}
		}

		// check se la querySQL è stata composta correttamente
		if query != "" {
			rows, err := models.DB.Query(query)
			if err != nil {
				log.Fatalf("Query: %v", err)
			}
			defer rows.Close()

			// counter per il conteggio dei record
			counter := 0

			// testing
			//log.Println(query)

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
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "errore formulazione richiesta",
			})
		}
	}

	return gin.HandlerFunc(fn)
}

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
