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

	// endpoint: /andamento/nazionale
	router.GET("/andamento/nazionale", controllers.NationalTrend)
	// endpoint: /andamento/nazionale/:bydate
	router.GET("/andamento/nazionale/:bydate", controllers.NationalTrendByDate)

	// endpoint: /utenti
	router.GET("/utenti", controllers.GetAllUsers)
	// endpoint: /utenti/:byusername
	router.GET("/utenti/:byusername", controllers.GetUserByUsername)

	router.Run(":" + port)
}
