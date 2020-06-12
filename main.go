package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"pdgt-covid/controllers"
	"pdgt-covid/models"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// db
	models.ConnectDatabase()

	router := gin.New()
	router.Use(gin.Logger())

	// endpoint iniziale
	router.GET("/", models.HandleWelcome)

	// endpoint andamenti
	router.GET("/andamento", controllers.HandleAndamento)
	router.GET("/andamento/nazionale", controllers.NationalTrend)
	router.GET("/andamento/nazionale/:bydate", controllers.NationalTrendByDate)

	// endpoint utenti
	router.GET("/utenti", controllers.GetAllUsers)
	router.GET("/utenti/:byusername", controllers.GetUserByUsername)
	router.POST("/utenti/registrazione", controllers.UserSignup)
	router.DELETE("/utenti/:byusername", controllers.UserDelete)

	router.Run(":" + port)
}
