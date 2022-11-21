package main

import (
	"log"
	"net/http"

	"github.com/anazworth/sleepStats/initializers"
	"github.com/anazworth/sleepStats/models"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
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

	r.POST("/api/v1/response", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		var responsebody struct {
			Response bool `json:"response"`
			Age      int  `json:"age"`
		}

		c.BindJSON(&responsebody)

		response := models.UserResponse{
			Response: responsebody.Response,
			Age:      responsebody.Age,
		}

		initializers.DB.Create(&response)

		c.JSON(200, gin.H{
			"message": response,
		})
	})

	r.GET("/api/v1/allResponses", func(c *gin.Context) {
		var responses []models.UserResponse

		if err := initializers.DB.Find(&responses).Error; err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, responses)
	})

	r.GET("/api/v1/summary", func(c *gin.Context) {
		var responses []models.UserResponse

		if err := initializers.DB.Find(&responses).Error; err != nil {
			log.Println(err)
		}

		var summary struct {
			Total      int     `json:"total"`
			TotalAdult int     `json:"totalAdult"`
			Yes        int     `json:"yes"`
			No         int     `json:"no"`
			NA         int     `json:"na"`
			PercentYes float64 `json:"percentYes"`
			PercentNo  float64 `json:"percentNo"`
		}

		summary.Total = len(responses)
		for _, response := range responses {
			if response.Response && response.Age > 17 {
				summary.Yes++
			} else if !response.Response && response.Age > 17 {
				summary.No++
			} else {
				summary.NA++
			}
		}
		summary.TotalAdult = summary.Total - summary.NA

		// Calculate percentages of valid 'yes/no' answers, rounding to 2 decimal places
		summary.PercentYes = float64(summary.Yes) / float64(summary.TotalAdult) * 100
		summary.PercentYes = float64(int(summary.PercentYes*100)) / 100
		summary.PercentNo = float64(summary.No) / float64(summary.TotalAdult) * 100
		summary.PercentNo = float64(int(summary.PercentNo*100)) / 100

		c.JSON(http.StatusOK, summary)
	})

	r.Run() // listen and serve on 8080
}
