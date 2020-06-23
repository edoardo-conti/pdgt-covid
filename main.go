package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"pdgt-covid/controllers"
	"pdgt-covid/middlewares"
	"pdgt-covid/models"

	_ "github.com/lib/pq"
)

func main() {
	// verifica della variabile d'ambiente che specifica la porta utilizzata per l'erogazione dal webservice
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// connessione al database Heroku Postgres
	models.ConnectDatabase()

	// inizializzazione router con tecnologia gin
	router := gin.New()
	router.Use(gin.Logger())

	// abilitazione richieste cors, per future modifiche fare riferimento a: https://github.com/gin-contrib/cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// endpoint iniziali
	router.GET("/api", models.HandleWelcome)
	router.GET("/api/trend", models.HandleAndamento)
	// endpoint trend nazionali
	router.GET("/api/trend/nazionale", controllers.NationalTrend)
	router.GET("/api/trend/nazionale/data/:bydate", controllers.NationalTrendByDate)
	router.GET("/api/trend/nazionale/picco", controllers.NationalTrendByPicco)
	router.POST("/api/trend/nazionale", middlewares.AuthMiddleware(), controllers.AddNationalTrend)
	router.DELETE("/api/trend/nazionale/data/:bydate", middlewares.AuthMiddleware(), controllers.DeleteNationalTrend)
	router.PATCH("/api/trend/nazionale/data/:bydate", middlewares.AuthMiddleware(), controllers.PatchNationalTrend)
	// endpoint trend regionali
	router.GET("/api/trend/regionale", controllers.RegionalTrendHandler("/"))
	router.GET("/api/trend/regionale/data/:bydata", controllers.RegionalTrendHandler("/data/:"))
	router.GET("/api/trend/regionale/regione/:byregid", controllers.RegionalTrendHandler("/regione/:"))
	router.GET("/api/trend/regionale/picco/", controllers.RegionalTrendHandler("/picco"))
	router.GET("/api/trend/regionale/picco/:byregid", controllers.RegionalTrendHandler("/picco/:"))
	// endpoint utenti
	router.GET("/api/utenti", middlewares.AuthMiddleware(), controllers.GetAllUsers)
	router.GET("/api/utenti/:byusername", middlewares.AuthMiddleware(), controllers.GetUserByUsername)
	router.POST("/api/utenti/signup", middlewares.AuthMiddleware(), controllers.UserSignup)
	router.POST("/api/utenti/login", controllers.UserSignin)
	router.DELETE("/api/utenti/:byusername", middlewares.AuthMiddleware(), controllers.UserDelete)

	// endpoint 404
	router.NoRoute(models.HandleNoRoute)

	// webservice...
	router.Run(":" + port)
}
