package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"pdgt-covid/controllers"
	"pdgt-covid/models"
)

type response1 struct {
	Status   int    `json:"status"`
	Messagge string `json:"message"`
}
type response2 struct {
	Status   int      `json:"status"`
	Messagge []string `json:"message"`
}

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
		str := `{"status": 200, "message": "Benvenuto su PDGT-COVID! [Developed by Edoardo C.]"}`
		res := response1{}
		json.Unmarshal([]byte(str), &res)

		c.JSON(http.StatusOK, res)
	})

	router.GET("/andamento", func(c *gin.Context) {
		str := `{"status": 200, "message": ["/andamento/nazionale", "/andamento/regionale"]}`
		res := response2{}
		json.Unmarshal([]byte(str), &res)

		c.JSON(http.StatusOK, res)
	})

	// endpoint: /nazione
	router.GET("/andamento/nazionale", controllers.NationalTrend)
	// endpoint: /nazione/:bydate (todo)
	router.GET("/andamento/nazionale/:bydate", controllers.NationalTrendByDate)

	router.Run(":" + port)
}
