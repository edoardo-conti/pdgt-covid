package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"pdgt-covid/middlewares"
	"pdgt-covid/models"
)

// GetAllUsers ottenere l'intera lista di utenti registrati nel sistema
func GetAllUsers(c *gin.Context) {
	// query SQL
	rows, err := models.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	/*
	 * gestione degli avatar tramite servizio esterno che genera
	 * l'immagine a partire dall'iniziale dell'username
	 */
	avatarURLBase := "https://avatars.dicebear.com/api/initials/"
	avatarURL := ""

	counter := 0

	// array di strutture utenti per raccogliere i risultati della query
	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.Username, &u.Password, &u.IsAdmin)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		// generazione URL avatar
		avatarURL = avatarURLBase + string([]rune(u.Username)[0]) + ".svg"

		users = append(users, models.User{Username: u.Username, Password: u.Password, IsAdmin: u.IsAdmin, Avatar: avatarURL})

		counter++
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"count":  counter,
		"data":   users,
	})
}

//GetUserByUsername ricerca di un utente nel database per username
func GetUserByUsername(c *gin.Context) {
	// si ricava il nome utente dal path dell'url RESTfull (:byusername)
	usrname := c.Params.ByName("byusername")

	// controllo veloce che il nome fornito non sia vuoto
	if usrname != "" {
		var u models.User
		// modello aggiuntivo utilizzato per l'integrazione dell'url dell'avatar
		var uc models.User

		// query SQL
		row := models.DB.QueryRow("SELECT * FROM users WHERE username=$1", usrname)
		switch err := row.Scan(&u.Username, &u.Password, &u.IsAdmin); err {
		case sql.ErrNoRows:
			// nessun record risultante
			c.JSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Utente richiesto non presente nel database.",
			})
		case nil:
			// generazione URL avatar
			avatarURL := "https://avatars.dicebear.com/api/initials/" + string([]rune(u.Username)[0]) + ".svg"
			// generazione struttura utente da restituire
			uc = models.User{Username: u.Username, Password: u.Password, IsAdmin: u.IsAdmin, Avatar: avatarURL}

			c.JSON(http.StatusOK, gin.H{
				"status": 200,
				"data":   uc,
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Errore nel risolvere la richiesta, riprova più tardi.",
			})
		}
	}
}

// UserSignup metodo utile alla registrazione di un utente nel sistema
func UserSignup(c *gin.Context) {
	// variabile per la validazione dell'input
	var newUserInput models.User
	// ShouldBindJSON si occupa di effettuare il binding tra richiesta POST e struttura dati fornita
	if err := c.ShouldBindJSON(&newUserInput); err == nil {
		// trim: eliminazione di spazi vuoti all'inizio e alla fine delle stringhe
		newUserInput.Username = strings.TrimSpace(newUserInput.Username)
		newUserInput.Password = strings.TrimSpace(newUserInput.Password)

		// verifica che i campi risultanti non siano vuoti
		if newUserInput.Username != "" && newUserInput.Password != "" {
			// generazione hash bcrypt della password fornita
			hashedpassword, err := bcrypt.GenerateFromPassword([]byte(newUserInput.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("Error hashing password: %q", err)
			}
			newUserInput.Password = string(hashedpassword)

			/*
			 * prima di procedere con l'inserimento dell'utente nel database si
			 * verifica che non esista già un utente con lo stesso username
			 */
			var u models.User
			row := models.DB.QueryRow("SELECT * FROM users WHERE username=$1", newUserInput.Username)
			switch err := row.Scan(&u.Username, &u.Password); err {
			case sql.ErrNoRows:
				// utente non presente, si può continuare con la registrazione
				_, err = models.DB.Exec("INSERT INTO users (username, password, isadmin) VALUES ($1, $2, $3);", newUserInput.Username, newUserInput.Password, newUserInput.IsAdmin)
				if err != nil {
					panic(err)
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status":  200,
						"message": "Utente registrato con successo.",
						"info":    "Per visualizzare: /utenti/" + newUserInput.Username,
					})
				}
			default:
				// utente già registrato nel database con l'username fornito
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "Errore: utente " + u.Username + " già registrato nel database.",
				})
			}
		} else {
			// uno o entrambi i campi sono vuoti
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  406,
				"message": "Errore: richiesti entrambi i campi.",
			})
		}
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  422,
			"message": "Errore: formato richiesta POST non corretta.",
		})
	}
}

//UserSignin metodo utile per accedere ai servizi come utente registrato
func UserSignin(c *gin.Context) {
	var newUserLogin models.User

	// ShouldBindJSON si occupa di effettuare il binding tra richiesta POST e struttura dati fornita
	if err := c.ShouldBindJSON(&newUserLogin); err == nil {
		// verifica che i campi non siano vuoti
		if newUserLogin.Username != "" && newUserLogin.Password != "" {
			// verifica che l'utente esista nel database
			var u models.User
			row := models.DB.QueryRow("SELECT password, isadmin FROM users WHERE username=$1", newUserLogin.Username)
			switch err := row.Scan(&u.Password, &u.IsAdmin); err {
			case sql.ErrNoRows:
				/*
				 * per questioni di sicurezza pur essendo a conoscenza che l'username
				 * non esiste nel database si restituisce un messaggio d'errore generico
				 * per evitare possibili attacci bruteforce sul nome utente.
				 */
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  401,
					"message": "Errore: credenziali errate, si prega di riprovare.",
				})
			case nil:
				// utente presente nel db

				// hash della password registrata nel db
				hashedPassword := []byte(u.Password)
				// password clear fornita nella richiesta POST
				password := []byte(newUserLogin.Password)

				// verifica dell'hashing delle password sfruttando metodo CompareHashAndPassword()
				if bcrypt.CompareHashAndPassword(hashedPassword, password) == nil {
					// credenziali corrette, si procede con la creazione del token che verrà in seguito restituito
					token, err := middlewares.CreateToken(newUserLogin.Username, u.IsAdmin)
					if err != nil {
						// imprevisto nella generazione del token
						c.JSON(http.StatusUnprocessableEntity, gin.H{
							"status":  422,
							"message": "Errore: " + err.Error(),
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"status":  200,
							"message": "Utente " + newUserLogin.Username + " loggato con successo.",
							"token":   token,
						})
					}
				} else {
					// credenziali errate, messaggio restituito con http status 'Unauthorized'
					c.JSON(http.StatusUnauthorized, gin.H{
						"status":  401,
						"message": "Errore: credenziali errate, si prega di riprovare.",
					})
				}
			default:
				// errore default
				c.JSON(http.StatusBadGateway, gin.H{
					"status":  502,
					"message": "Errore: risultato inaspettato, si prega di riprovare più tardi.",
				})
			}
		} else {
			// uno o entrambi i campi sono vuoti
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  406,
				"message": "Errore: richiesti entrambi i campi.",
			})
		}
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  422,
			"message": "Errore: formato richiesta POST non corretta.",
		})
	}
}

//UserDelete metodo dedicato all'eliminazione di un utente dal sistema
func UserDelete(c *gin.Context) {
	// si ricava il nome utente dal path dell'url RESTfull (:byusername)
	usernameToDelete := strings.TrimSpace(c.Params.ByName("byusername"))

	// verifica che il parametro fornito non sia vuoto
	if usernameToDelete != "" {
		// verifica che l'utente con tale username esista nel db
		var u models.User
		row := models.DB.QueryRow("SELECT * FROM users WHERE username=$1", usernameToDelete)
		switch err := row.Scan(&u.Username, &u.Password, &u.IsAdmin); err {
		case sql.ErrNoRows:
			// utente non trovato
			c.JSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Utente " + usernameToDelete + " non registrato nel database.",
			})
		case nil:
			// l'utente esiste, si procede il delete
			res, err := models.DB.Exec("DELETE FROM users WHERE username=$1", usernameToDelete)
			if err == nil {
				count, err := res.RowsAffected()
				if err == nil {
					if count == 1 {
						c.JSON(http.StatusOK, gin.H{
							"status":  200,
							"message": "Utente " + usernameToDelete + " eliminato dal database con successo.",
						})
					}
				}
			} else {
				// gestire errore
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "Errore: si prega di riprovare più tardi.",
				})
			}
		default:
			// errore default
			c.JSON(http.StatusBadGateway, gin.H{
				"status":  502,
				"message": "Errore: risultato inaspettato, si prega di riprovare più tardi.",
			})
		}
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  406,
			"message": "Errore: campo vuoto non accettato.",
		})
	}
}
