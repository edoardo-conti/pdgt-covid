package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"pdgt-covid/models"
	"pdgt-covid/controllers"

	// importare: "github.com/edoardo-conti/pdgt-covid/models"
	// importare: "github.com/edoardo-conti/pdgt-covid/controllers"
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

	// endpoint: /
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Benvenuto su PDGT-COVID!",
		})
	})

	// endpoint: /nazione
	router.GET("/nazione", controllers.NationalTrend)

	// endpoint: /nazione/:bydate (todo)
	router.GET("/nazione/:bydate", controllers.NationalTrendByDate)

	router.Run(":" + port)
}
