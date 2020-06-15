package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"pdgt-covid/models"
)

// NationalTrend ...
func NationalTrend(c *gin.Context) {
	rows, err := models.DB.Query("SELECT * FROM nazione ORDER BY data ASC")
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

func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := models.DB.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if row exists '%s' %v", args, err)
	}
	return exists
}

//todo: controllare validazione post
func checkNewTrendInput(nti models.NationalTrendPOST) bool {
	check := false

	if nti.Data != "" &&
		nti.RicoveratiConSintomi >= 0 &&
		nti.TerapiaIntensiva >= 0 &&
		nti.TotaleOspedalizzati >= 0 &&
		nti.IsolamentoDomiciliare >= 0 &&
		nti.TotalePositivi >= 0 &&
		nti.NuoviPositivi >= 0 &&
		nti.DimessiGuariti >= 0 &&
		nti.Deceduti >= 0 &&
		nti.TotaleCasi >= 0 &&
		nti.Tamponi >= 0 &&
		nti.CasiTestati >= 0 {
		check = true
	}

	return check
}

//AddNationalTrend ...
func AddNationalTrend(c *gin.Context) {
	// Validate input
	var newTrendInput models.NationalTrendPOST
	if err := c.ShouldBindJSON(&newTrendInput); err == nil {
		// trim string fields
		newTrendInput.Data = strings.TrimSpace(newTrendInput.Data)

		// check valid fields
		//check := checkNewTrendInput(newTrendInput)
		dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
		//if check && newTrendInput.Data != "" && dateCheck.MatchString(newTrendInput.Data) {
		if newTrendInput.Data != "" && dateCheck.MatchString(newTrendInput.Data) {
			// check if trend already exist with the same date
			if rowExists("SELECT 1 FROM nazione WHERE data=$1", newTrendInput.Data) {
				// trend on that date already registered
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "trend in data " + newTrendInput.Data + " già registrato nel database",
					"info":    "/andamento/nazionale/data/" + newTrendInput.Data,
				})
			} else {
				// trend not found, can proceed
				_, err = models.DB.Exec("INSERT INTO nazione VALUES ($1, 'ITA', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);",
					newTrendInput.Data,
					newTrendInput.RicoveratiConSintomi,
					newTrendInput.TerapiaIntensiva,
					newTrendInput.TotaleOspedalizzati,
					newTrendInput.IsolamentoDomiciliare,
					newTrendInput.TotalePositivi,
					newTrendInput.VariazioneTotalePositivi,
					newTrendInput.NuoviPositivi,
					newTrendInput.DimessiGuariti,
					newTrendInput.Deceduti,
					newTrendInput.TotaleCasi,
					newTrendInput.Tamponi,
					newTrendInput.CasiTestati,
					"",
					"",
				)
				if err != nil {
					panic(err)
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status":  200,
						"message": "trend giornaliero registrato con successo",
						"info":    "per visualizzare: /andamento/nazionale/data/" + newTrendInput.Data,
					})
				}
			}
		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  406,
				"message": "richiesti tutti i campi nei rispettivi formati",
			})
		}
	} else {
		// (todo) try: c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "formato richiesta POST non corretta",
		})
	}
}

//DeleteNationalTrend ...
func DeleteNationalTrend(c *gin.Context) {
	// get parameter
	trendToDelete := strings.TrimSpace(c.Params.ByName("bydate"))

	// check valid field
	dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
	if trendToDelete != "" && dateCheck.MatchString(trendToDelete) {
		// check if trend (date) exist
		if rowExists("SELECT 1 FROM nazione WHERE data=$1", trendToDelete) {
			// trend on that date exist
			res, err := models.DB.Exec("DELETE FROM nazione WHERE data=$1", trendToDelete)
			if err == nil {
				count, err := res.RowsAffected()
				if err == nil {
					if count == 1 {
						c.JSON(http.StatusOK, gin.H{
							"status":  200,
							"message": "trend in data " + trendToDelete + " eliminato dal database con successo",
						})
					}
				}
			} else {
				// gestire errore (todo)
			}
		} else {
			// trend don't exist on that date
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "trend in data " + trendToDelete + " non presente nel database",
			})
		}
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  406,
			"message": "parametro non conforme",
		})
	}
}

//generateUpdateQuery utile a generare la query (stringa) per l'aggiornamento di un trend nazionale giornaliero
func generateUpdateQuery(ntu models.NationalTrendPATCH, dttu string) (query string, err error) {
	query = "UPDATE nazione SET"

	// marshal the struct
	fieldsToUpdate, err := json.Marshal(ntu)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(fieldsToUpdate))

	// iterate through the body of the post req.
	var v interface{}
	json.Unmarshal(fieldsToUpdate, &v)
	data := v.(map[string]interface{})

	// check body lenght
	counter := 0

	for k, v := range data {
		//fmt.Println(k, v)
		if k != "variazione_totale_positivi" && v.(float64) < 0 {
			return "", errors.New("valore campo negativo non permesso")
		}
		query += " " + fmt.Sprintf("%v", k) + "=" + fmt.Sprintf("%v", v) + ","

		counter++
	}

	if counter == 0 {
		return "", errors.New("body vuoto")
	}

	// get rid of last AND from the query
	query = strings.TrimRight(query, ",")

	// add the where clause
	query += " WHERE data='" + dttu + "';"

	return query, nil
}

//PatchNationalTrend ...
func PatchNationalTrend(c *gin.Context) {
	// get parameter
	dataTrendToUpdate := strings.TrimSpace(c.Params.ByName("bydate"))

	// check valid field
	dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
	if dataTrendToUpdate != "" && dateCheck.MatchString(dataTrendToUpdate) {
		// check if trend (date) exist
		if rowExists("SELECT 1 FROM nazione WHERE data=$1", dataTrendToUpdate) {
			// trend on that date exist

			// check BODY request
			var newTrendUpdate models.NationalTrendPATCH
			if err := c.ShouldBindJSON(&newTrendUpdate); err == nil {
				// logging struct with variable names
				// fmt.Printf("%+v\n", newTrendUpdate)

				upQuery, err1 := generateUpdateQuery(newTrendUpdate, dataTrendToUpdate)
				if err1 == nil {
					// query generata con successo!
					fmt.Println(upQuery)

					res, err2 := models.DB.Exec(upQuery)
					if err2 == nil {
						count, err3 := res.RowsAffected()
						if err3 == nil {
							if count == 1 {
								c.JSON(http.StatusOK, gin.H{
									"status":  200,
									"message": "trend in data " + dataTrendToUpdate + " aggiornato con successo",
									"info":    "/andamento/nazionale/data/" + dataTrendToUpdate,
								})
							}
						}
					} else {
						// gestire errore (todo)
					}
				} else {
					//fmt.Println("errore, uno dei valori è negativo")
					c.JSON(http.StatusNotAcceptable, gin.H{
						"status":  406,
						"message": "uno dei parametri dei campi non è conforme",
					})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "formato richiesta PATCH non corretta",
				})
			}
		} else {
			// trend don't exist on that date
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "trend in data " + dataTrendToUpdate + " non presente nel database",
			})
		}
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  406,
			"message": "url non conforme",
		})
	}
}
