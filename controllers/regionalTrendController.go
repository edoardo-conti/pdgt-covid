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
 * variabile utile a raggruppare la prima parte di una query SQL comune a diverse richieste ( “write less, do more" )
 * la query in questione genera e restituisce un oggetto json popolato con i campi interessati
 * raggruppato dalla data di rilevazione dell'andamento, ottenendo trend aggregati per data
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

//RegionalTrendHandler metodo dedicato alla gestione delle richieste GET riguardo al trend regionale
func RegionalTrendHandler(mode string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		/*
		 * multiplexing delle richieste GET suddiviso dal selettore di modalità 'mode'
		 * metodo sfruttato per servire con un'unica query di base tutti i GET degli endpoind trend regionale
		 */

		// preparazione query effettiva
		query := ""

		// selettore di modalità
		if mode == "/" {
			// query NON filtrata
			query = queryFirstPart + " GROUP BY regioni.data ORDER BY regioni.data ASC;"
		} else if mode == "/data/:" {
			// filtraggio per data
			data := c.Params.ByName("bydata")

			// check validità parametro (es. 2020-02-24)
			dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
			if data != "" && dateCheck.MatchString(data) {
				// se data in formato corretto si genera la query
				query = queryFirstPart + " WHERE data='" + data + "' GROUP BY regioni.data ORDER BY regioni.data ASC;"
			}
		} else if mode == "/regione/:" {
			// filtraggio per id regione
			regID := c.Params.ByName("byregid")

			// check validità parametro
			if regIDint, err := strconv.Atoi(regID); err == nil {
				if regIDint >= 0 && regIDint <= 22 {
					// generazione query
					query = queryFirstPart + " WHERE codice_regione='" + regID + "' GROUP BY regioni.data ORDER BY regioni.data ASC;"
				}
			}
		} else if mode == "/picco" {
			// filtraggio per picco
			// generazione query sfruttando funzione sql MAX() per il calcolo dei nuovi positivi totali, ergo picco richiesto
			query = queryFirstPart + " WHERE nuovi_positivi=(SELECT MAX(nuovi_positivi) FROM regioni) GROUP BY regioni.data ORDER BY regioni.data ASC;"
		} else if mode == "/picco/:" {
			// filtraggio per picco di regione
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
			// inoltro query
			rows, err := models.DB.Query(query)
			if err != nil {
				log.Fatalf("Query: %v", err)
			}
			defer rows.Close()

			// counter per il conteggio dei record
			counter := 0

			// struttura dedicata alla "collezzione" dei trend regionali
			var regioniCollect []models.RegionalTrendCollect
			for rows.Next() {
				var r models.RegionalTrendCollect

				err = rows.Scan(&r.Data, &r.Info)
				if err != nil {
					log.Fatalf("scan error: %v", err)
				}

				// si popola la struttura collezzione dei record restituiti dalla query
				regioniCollect = append(regioniCollect, models.RegionalTrendCollect{Data: r.Data, Info: r.Info})

				counter++
			}

			// verifica dei record ottenuti
			if counter > 0 {
				c.JSON(http.StatusOK, gin.H{
					"status": 200,
					"count":  counter,
					"data":   regioniCollect,
				})
			} else {
				// zero record restituiti
				c.JSON(http.StatusNotFound, gin.H{
					"status":  404,
					"message": "Errore: la richiesta non ha prodotto risultati.",
				})
			}
		} else {
			// query risultante nulla
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Errore: formulazione richiesta non corretta.",
			})
		}
	}

	return gin.HandlerFunc(fn)
}
