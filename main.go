package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"pdgt-covid/controllers"
	"pdgt-covid/middlewares"
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
	router.GET("/utenti", middlewares.AuthMiddleware(), controllers.GetAllUsers)
	router.GET("/utenti/:byusername", middlewares.AuthMiddleware(), controllers.GetUserByUsername)
	router.POST("/utenti/signup", controllers.UserSignup)
	router.POST("/utenti/login", controllers.UserSignin)
	router.DELETE("/utenti/:byusername", middlewares.AuthMiddleware(), controllers.UserDelete)

	// endpoint 404
	router.NoRoute(models.HandleNoRoute)

	router.Run(":" + port)
}
