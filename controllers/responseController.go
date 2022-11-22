package controllers

import (
	"log"
	"net/http"

	"github.com/anazworth/sleepStats/initializers"
	"github.com/anazworth/sleepStats/models"
	"github.com/gin-gonic/gin"
)

func GetAllResponses(c *gin.Context) {
	var responses []models.UserResponse

	if err := initializers.DB.Find(&responses).Error; err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, responses)
}

func CreateResponse(c *gin.Context) {
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
}
