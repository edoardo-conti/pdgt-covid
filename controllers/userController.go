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

// GetAllUsers ...
// @returns var -> descrizione param and 2nd var -> desc
func GetAllUsers(c *gin.Context) {
	rows, err := models.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	var users []models.User
	counter := 0

	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.Username, &u.Password)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		users = append(users, models.User{u.Username, u.Password})

		counter++
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"count":  counter,
		"data":   users,
	})
}

//GetUserByUsername ...
func GetUserByUsername(c *gin.Context) {
	// get parameter
	usrname := c.Params.ByName("byusername")

	if usrname != "" {
		// controllo validità del parametro (todo)
		var u models.User

		row := models.DB.QueryRow("SELECT * FROM users WHERE username=$1", usrname)
		switch err := row.Scan(&u.Username, &u.Password); err {
		case sql.ErrNoRows:
			c.JSON(http.StatusOK, gin.H{
				"status":  200,
				"message": "utente richiesto non disponibile",
			})
		case nil:
			c.JSON(http.StatusOK, gin.H{
				"status": 200,
				"data":   u,
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "formato richiesta non corretto",
			})
		}
	}
}

//UserSignup ...
func UserSignup(c *gin.Context) {
	// Validate input
	var newUserInput models.User
	if err := c.ShouldBindJSON(&newUserInput); err == nil {
		// trim fields
		newUserInput.Username = strings.TrimSpace(newUserInput.Username)
		newUserInput.Password = strings.TrimSpace(newUserInput.Password)

		// check valid fields
		if newUserInput.Username != "" && newUserInput.Password != "" {
			// hash password
			hashedpassword, err := bcrypt.GenerateFromPassword([]byte(newUserInput.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("Error hashing password: %q", err)
			}
			newUserInput.Password = string(hashedpassword)

			// check if username already exist
			var u models.User
			row := models.DB.QueryRow("SELECT * FROM users WHERE username=$1", newUserInput.Username)
			switch err := row.Scan(&u.Username, &u.Password); err {
			case sql.ErrNoRows:
				// users not found, can proceed
				_, err = models.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2);", newUserInput.Username, newUserInput.Password)
				if err != nil {
					panic(err)
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status":  200,
						"message": "utente registrato con successo",
						"info":    "per visualizzare: /utenti/" + newUserInput.Username,
					})
				}
			case nil:
				// user with that username already registered
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "username già registrato nel database",
				})
			default:
				// gestire errore (todo)
			}

		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  406,
				"message": "richiesti entrambi i campi",
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

//UserSignin ...
func UserSignin(c *gin.Context) {
	// Validate input
	var newUserLogin models.User
	if err := c.ShouldBindJSON(&newUserLogin); err == nil {
		// check valid fields
		if newUserLogin.Username != "" && newUserLogin.Password != "" {
			// check if user exist
			var u models.User
			row := models.DB.QueryRow("SELECT password FROM users WHERE username=$1", newUserLogin.Username)
			switch err := row.Scan(&u.Password); err {
			case sql.ErrNoRows:
				// (todo) da correggere per sicurezza: users not found
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  400,
					"message": "utente non trovato",
				})
			case nil:
				// user exist

				// Get the hashed password from the saved document
				hashedPassword := []byte(u.Password)
				// Get the password provided in the request.body
				password := []byte(newUserLogin.Password)

				// check user password
				if bcrypt.CompareHashAndPassword(hashedPassword, password) == nil {
					// credenziali corrette, procedo
					token, err := middlewares.CreateToken(newUserLogin.Username)
					if err != nil {
						c.JSON(http.StatusUnprocessableEntity, gin.H{
							"status":  422,
							"message": err.Error(),
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"status":  200,
							"message": "utente loggato con successo",
							"token":   token,
						})
					}
				} else {
					// credenziali errate
					c.JSON(http.StatusUnauthorized, gin.H{
						"status":  401,
						"message": "credenziali errate",
					})
				}
			default:
				// gestire errore (todo)
				log.Println("errore")
			}
		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  406,
				"message": "richiesti entrambi i campi",
			})
		}
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  422,
			"message": "formato richiesta POST non corretta",
		})
	}
}

//UserDelete ...
func UserDelete(c *gin.Context) {
	// get parameter
	usernameToDelete := strings.TrimSpace(c.Params.ByName("byusername"))

	// check valid field
	if usernameToDelete != "" {
		// check if username exist
		var u models.User
		row := models.DB.QueryRow("SELECT * FROM users WHERE username=$1", usernameToDelete)
		switch err := row.Scan(&u.Username, &u.Password); err {
		case sql.ErrNoRows:
			// users not found
			c.JSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "username non registrato nel database",
			})
		case nil:
			// user exist
			res, err := models.DB.Exec("DELETE FROM users WHERE username=$1", usernameToDelete)
			if err == nil {
				count, err := res.RowsAffected()
				if err == nil {
					if count == 1 {
						c.JSON(http.StatusOK, gin.H{
							"status":  200,
							"message": "utente eliminato dal database con successo",
						})
					}
				}
			} else {
				// gestire errore (todo)
			}
		default:
			// gestire errore (todo)
		}

	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  406,
			"message": "campo vuoto non accettato",
		})
	}
}
