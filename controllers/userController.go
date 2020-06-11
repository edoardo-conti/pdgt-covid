package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"pdgt-covid/models"

	"github.com/gin-gonic/gin"
)

// GetAllUsers ...
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
