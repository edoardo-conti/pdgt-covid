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
	router.GET("/andamento", models.HandleAndamento)
	router.GET("/andamento/nazionale", controllers.NationalTrend)
	router.GET("/andamento/nazionale/data/:bydate", controllers.NationalTrendByDate)
	router.GET("/andamento/nazionale/picco", controllers.NationalTrendByPicco)
	router.POST("/andamento/nazionale", controllers.AddNationalTrend)
	router.DELETE("/andamento/nazionale/data/:bydate", controllers.DeleteNationalTrend)
	router.PATCH("/andamento/nazionale/data/:bydate", controllers.PatchNationalTrend)

	router.GET("/andamento/regionale", controllers.RegionalTrendHandler(1))
	router.GET("/andamento/regionale/data/:bydata", controllers.RegionalTrendHandler(2))
	router.GET("/andamento/regionale/regione/:byregid", controllers.RegionalTrendHandler(3))
	router.GET("/andamento/regionale/picco/", controllers.RegionalTrendHandler(4))
	router.GET("/andamento/regionale/picco/:byregid", controllers.RegionalTrendHandler(5))

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
