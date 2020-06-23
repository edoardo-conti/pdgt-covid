package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"pdgt-covid/models"
)

// NationalTrend gestione metodo GET riguardo al trend nazionale
func NationalTrend(c *gin.Context) {
	// query SQL
	rows, err := models.DB.Query("SELECT * FROM nazione ORDER BY data ASC")
	if err != nil {
		log.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	// struttura dati dedicata allo storage dei record rilevati
	var nazioni []models.NationalTrend

	// counter per il conteggio dei record
	counter := 0

	// si cicla per ogni record ottenuto dalla query precedente
	for rows.Next() {
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

		// generazione oggetto di risposta
		nazioni = append(nazioni, models.NationalTrend{
			Data:                     r.Data,
			Stato:                    r.Stato,
			RicoveratiConSintomi:     r.RicoveratiConSintomi,
			TerapiaIntensiva:         r.TerapiaIntensiva,
			TotaleOspedalizzati:      r.TotaleOspedalizzati,
			IsolamentoDomiciliare:    r.IsolamentoDomiciliare,
			TotalePositivi:           r.TotalePositivi,
			VariazioneTotalePositivi: r.VariazioneTotalePositivi,
			NuoviPositivi:            r.NuoviPositivi,
			DimessiGuariti:           r.DimessiGuariti,
			Deceduti:                 r.Deceduti,
			TotaleCasi:               r.TotaleCasi,
			Tamponi:                  r.Tamponi,
			CasiTestati:              r.CasiTestati,
			NoteIT:                   r.NoteIT,
			NoteEN:                   r.NoteEN})

		// incremento del counter
		counter++
	}

	// restituzione trend nazionali rilevati dal database
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"count":  counter,
		"data":   nazioni,
	})
}

// NationalTrendByDate metodo dedicato al filtraggio dei trend nazionali per data
func NationalTrendByDate(c *gin.Context) {
	// parametro data ottenuto da path url
	date := c.Params.ByName("bydate")

	// verifica che il parametro non sia vuoto
	if date != "" {
		// controllo validità formato data (es. 2020-04-30)
		dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
		if dateCheck.MatchString(date) {
			var r models.NationalTrend

			row := models.DB.QueryRow("SELECT * FROM nazione WHERE data=$1 LIMIT 1", date)
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
				// nessun record risultate
				c.JSON(http.StatusNotFound, gin.H{
					"status":  404,
					"message": "Trend nazionale in data " + date + " non disponibile.",
				})
			case nil:
				// trend risultante
				c.JSON(http.StatusOK, gin.H{
					"status": 200,
					"data":   r,
				})
			default:
				// errore default
				c.JSON(http.StatusBadGateway, gin.H{
					"status":  502,
					"message": "Errore: risultato inaspettato, si prega di riprovare più tardi.",
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Errore: formato data fornita non corretto.",
				"format":  "es. 2020-04-12",
			})
		}
	}
}

// NationalTrendByPicco filtraggio trend nazionale per picco di nuovi positivi
func NationalTrendByPicco(c *gin.Context) {
	var r models.NationalTrend

	// query SQL per la ricerca del massimo valore di nuovi_positivi (picco) tra i record registrati
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
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "Errore: al momento non è possibile trovare il record picco.",
		})
	case nil:
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"data":   r,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Errore: si prega di riprovare più tardi.",
		})
	}
}

// rowExists metodo dedicato alla verifica rapida dell'esistenza di un record nel database
func rowExists(query string, args ...interface{}) bool {
	var exists bool

	// si sfrutta funzione exists() SQL
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := models.DB.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Errore: verifica della query fallita: '%s' - %v", args, err)
	}

	return exists
}

//checkAddTrendFields metodo dedicato alla verifica dei parametri della richiesta POST di un trend nazionale
func checkAddTrendFields(ntp models.NationalTrendPOST) bool {
	ret := false

	// conversione dei parametri da stringa ad intero
	TerapiaIntensiva, eti := strconv.Atoi(ntp.TerapiaIntensiva)
	RicoveratiConSintomi, ers := strconv.Atoi(ntp.RicoveratiConSintomi)
	TotaleOspedalizzati, eto := strconv.Atoi(ntp.TotaleOspedalizzati)
	IsolamentoDomiciliare, eid := strconv.Atoi(ntp.IsolamentoDomiciliare)
	TotalePositivi, etp := strconv.Atoi(ntp.TotalePositivi)
	NuoviPositivi, enp := strconv.Atoi(ntp.NuoviPositivi)
	DimessiGuariti, edg := strconv.Atoi(ntp.DimessiGuariti)
	Deceduti, ed := strconv.Atoi(ntp.Deceduti)
	TotaleCasi, etc := strconv.Atoi(ntp.TotaleCasi)
	Tamponi, et := strconv.Atoi(ntp.Tamponi)
	CasiTestati, ect := strconv.Atoi(ntp.CasiTestati)

	// verifica che le stringhe convertite contengano numeri (isNaN)
	if eti == nil &&
		ers == nil &&
		eto == nil &&
		eid == nil &&
		etp == nil &&
		enp == nil &&
		edg == nil &&
		ed == nil &&
		etc == nil &&
		et == nil &&
		ect == nil {
		// variazione_totale_positivi unico campo che può essere negativo
		// si controllano tutti i restanti per verificare che siano positivi
		if RicoveratiConSintomi >= 0 &&
			TerapiaIntensiva >= 0 &&
			TotaleOspedalizzati >= 0 &&
			IsolamentoDomiciliare >= 0 &&
			TotalePositivi >= 0 &&
			NuoviPositivi >= 0 &&
			DimessiGuariti >= 0 &&
			Deceduti >= 0 &&
			TotaleCasi >= 0 &&
			Tamponi >= 0 &&
			CasiTestati >= 0 {
			ret = true
		}
	}

	return ret
}

//AddNationalTrend metodo dedicato alla gestione delle richieste POST per i trend nazionali
func AddNationalTrend(c *gin.Context) {
	var newTrendInput models.NationalTrendPOST

	/*
	 * bind tra sfruttura dati e json contenuti nel body della richiesta POST
	 * in ogni campo della struttura dati 'NationalTrendPOST' è presente il tag 'binding:"required"'
	 * per tanto ognuno di questi è necessario per concludere con successo un metodo POST
	 */
	if err := c.ShouldBindJSON(&newTrendInput); err == nil {
		// trim del campo data
		newTrendInput.Data = strings.TrimSpace(newTrendInput.Data)

		// check dei campo data (es. 2020-02-24) e verifica dei campi tramite metodo checkAddTrendFields()
		dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
		if newTrendInput.Data != "" && dateCheck.MatchString(newTrendInput.Data) && checkAddTrendFields(newTrendInput) {
			// si verifica che non esista già un record con la stessa data sfruttando il metodo rowExists()
			if rowExists("SELECT 1 FROM nazione WHERE data=$1", newTrendInput.Data) {
				// trend in data richiesta già disponibile
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "Trend in data " + newTrendInput.Data + " già registrato nel database.",
					"info":    "/api/trend/nazionale/data/" + newTrendInput.Data,
				})
			} else {
				// trend non trovato in data richiesta, si può procedere
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
					nil,
					nil,
				)
				if err != nil {
					// errore default
					c.JSON(http.StatusBadGateway, gin.H{
						"status":  502,
						"message": "Errore: risultato inaspettato, si prega di riprovare più tardi.",
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status":  200,
						"message": "Trend giornaliero nazionale registrato con successo.",
						"info":    "/api/trend/nazionale/data/" + newTrendInput.Data,
					})
				}
			}
		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  406,
				"message": "Errore: uno o più parametri forniti non sono conformi al formato richiesto.",
			})
		}
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  422,
			"message": "Errore: formato richiesta POST non corretta (campi omessi o formati non corretti).",
		})
	}
}

//DeleteNationalTrend metodo dedicato all'eliminazione di un trend nazionale rilevato in data fornita
func DeleteNationalTrend(c *gin.Context) {
	// si ricava la data del record da eliminare
	trendToDelete := strings.TrimSpace(c.Params.ByName("bydate"))

	// verifica validità parametro
	dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
	if trendToDelete != "" && dateCheck.MatchString(trendToDelete) {
		// si verifica se il trend in data selezionata esista
		if rowExists("SELECT 1 FROM nazione WHERE data=$1", trendToDelete) {
			// in caso positivo si procede con l'eliminazione via query SQL
			res, err := models.DB.Exec("DELETE FROM nazione WHERE data=$1", trendToDelete)
			if err == nil {
				count, err := res.RowsAffected()
				if err == nil {
					if count == 1 {
						c.JSON(http.StatusOK, gin.H{
							"status":  200,
							"message": "Trend in data " + trendToDelete + " eliminato dal database con successo.",
						})
					}
				}
			} else {
				c.JSON(http.StatusBadGateway, gin.H{
					"status":  502,
					"message": "Errore: risultato inaspettato, si prega di riprovare più tardi.",
				})
			}
		} else {
			// trend non disponibile in data richiesta
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Errore: trend in data " + trendToDelete + " non presente nel database.",
			})
		}
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  406,
			"message": "Errore: parametro fornito non corretto.",
		})
	}
}

//generateUpdateQuery utile a generare dinamicamente la query dell'aggiornamento (patch) di un trend nazionale giornaliero
func generateUpdateQuery(ntu models.NationalTrendPATCH, dttu string) (query string, err error) {
	var v interface{}

	// inizializzazione stringa query SQL
	query = "UPDATE nazione SET"

	// marshal ed unmarshal della struttura (utile per sfruttare tag 'omitempty')
	fieldsToUpdate, err := json.Marshal(ntu)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(fieldsToUpdate, &v)

	// conteggio campi nel body
	counter := 0

	// scorrimento dei campi nel body della richiesta POST
	data := v.(map[string]interface{})
	for k, v := range data {
		if k != "variazione_totale_positivi" && v.(float64) < 0 {
			return "", errors.New("valore campo negativo non permesso")
		}

		// composizione query SQL aggiungendo campo e valore direttamente da richiesta body
		query += " " + fmt.Sprintf("%v", k) + "=" + fmt.Sprintf("%v", v) + ","

		counter++
	}

	// verifica del numero di iterazioni
	if counter == 0 {
		return "", errors.New("body vuoto")
	}

	// rimozione dell'ultima virgola in eccesso dalla query generata dinamicamente
	query = strings.TrimRight(query, ",")
	// aggiunta della clausola SQl WHERE
	query += " WHERE data='" + dttu + "';"

	return query, nil
}

//PatchNationalTrend metodo dedicato all'aggiornamento dei campi di un record pre-esistente nel database
func PatchNationalTrend(c *gin.Context) {
	// si ricava il parametro data
	dataTrendToUpdate := strings.TrimSpace(c.Params.ByName("bydate"))

	// verifica validità formato data
	dateCheck := regexp.MustCompile("((19|20)\\d\\d)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])")
	if dataTrendToUpdate != "" && dateCheck.MatchString(dataTrendToUpdate) {
		// si verifica che il trend esista in data fornita
		if rowExists("SELECT 1 FROM nazione WHERE data=$1", dataTrendToUpdate) {
			// trend esiste, si procede
			// verifica del body della richiesta PATCH sfruttando ShouldBindJSON()
			var newTrendUpdate models.NationalTrendPATCH
			if err := c.ShouldBindJSON(&newTrendUpdate); err == nil {
				// generazione dinamica della query SQL di aggiornamento
				upQuery, err1 := generateUpdateQuery(newTrendUpdate, dataTrendToUpdate)
				if err1 == nil && upQuery != "" {
					// query generata con successo, si procedere inoltrando la richiesta
					res, err2 := models.DB.Exec(upQuery)
					if err2 == nil {
						count, err3 := res.RowsAffected()
						if err3 == nil {
							if count == 1 {
								c.JSON(http.StatusOK, gin.H{
									"status":  200,
									"message": "Trend in data " + dataTrendToUpdate + " aggiornato con successo.",
									"info":    "/api/trend/nazionale/data/" + dataTrendToUpdate,
								})
							}
						}
					} else {
						c.JSON(http.StatusBadGateway, gin.H{
							"status":  502,
							"message": "Errore: risultato inaspettato, si prega di riprovare più tardi.",
						})
					}
				} else {
					c.JSON(http.StatusNotAcceptable, gin.H{
						"status":  406,
						"message": "Errore: uno dei parametri forniti non risulta essere corretto.",
					})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "Errore: formato richiesta PATCH non corretta.",
				})
			}
		} else {
			// trend in data selezionata non presente
			c.JSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Errore: trend in data " + dataTrendToUpdate + " non presente nel database.",
			})
		}
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  406,
			"message": "Errore: formato data fornita non corretto.",
		})
	}
}
