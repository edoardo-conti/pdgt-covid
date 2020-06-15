package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"pdgt-covid/models"
)

// NationalTrend ...
func NationalTrend(c *gin.Context) {
	rows, err := models.DB.Query("SELECT * FROM nazione")
	if err != nil {
		log.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	//got := []nazione{}

	var nazioni []models.NationalTrend
	// counter per il conteggio dei record
	counter := 0

	for rows.Next() {
		// c := new(Course)
		var r models.NationalTrend
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
		nazioni = append(nazioni, models.NationalTrend{r.Data, r.Stato, r.RicoveratiConSintomi, r.TerapiaIntensiva, r.TotaleOspedalizzati, r.IsolamentoDomiciliare, r.TotalePositivi, r.VariazioneTotalePositivi, r.NuoviPositivi, r.DimessiGuariti, r.Deceduti, r.TotaleCasi, r.Tamponi, r.CasiTestati, r.NoteIT, r.NoteEN})

		// incremento del counter
		counter++

		//got = append(got, r)
	}

	//log.Println(got)
	//nazioniBytes, _ := json.Marshal(&nazioni)

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"count":  counter,
		"data":   nazioni,
	})
}

// NationalTrendByDate ...
func NationalTrendByDate(c *gin.Context) {
	// get parameter
	date := c.Params.ByName("bydate")

	//fmt.Println("log: %s", date)

	if date != "" {
		// controllo validità del parametro (es. 2020-04-30)
		dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
		if dateCheck.MatchString(date) {
			var r models.NationalTrend

			row := models.DB.QueryRow("SELECT * FROM nazione WHERE data=$1", date)
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
				c.JSON(http.StatusOK, gin.H{
					"status":  200,
					"message": "trend data richiesta non disponibile",
				})
			case nil:
				c.JSON(http.StatusOK, gin.H{
					"status": 200,
					"data":   r,
				})
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "formato data non corretto",
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "formato data non corretto",
				"format":  "es. 2020-04-12",
			})
		}
	}
}

// NationalTrendByPicco ...
func NationalTrendByPicco(c *gin.Context) {
	var r models.NationalTrend

	row := models.DB.QueryRow("SELECT * FROM nazione WHERE nuovi_positivi=(select max(nuovi_positivi) from nazione)")
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
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "trend data richiesta non disponibile",
		})
	case nil:
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"data":   r,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "formato data non corretto",
		})
	}
}
