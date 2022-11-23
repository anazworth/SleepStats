package main

import (
	"github.com/anazworth/sleepStats/controllers"
	"github.com/anazworth/sleepStats/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	// initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.MigrateDB()
}

func main() {

	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "alive",
		})
	})

	r.POST("/api/v1/response", controllers.CreateResponse)

	r.GET("/api/v1/allResponses", controllers.GetAllResponses)

	r.GET("/api/v1/summary", controllers.GetSummary)

	r.GET("/api/v1/dataInterpretation", controllers.InterperetData)

	r.Run() // listen and serve on 8080
}
