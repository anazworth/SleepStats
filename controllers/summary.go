package controllers

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/anazworth/sleepStats/initializers"
	"github.com/anazworth/sleepStats/models"
	"github.com/gin-gonic/gin"
	"github.com/montanaflynn/stats"
)

func GetSummary(c *gin.Context) {
	var responses []models.UserResponse

	var Summary struct {
		Total      int     `json:"total"`
		TotalAdult int     `json:"totalAdult"`
		Yes        int     `json:"yes"`
		No         int     `json:"no"`
		NA         int     `json:"na"`
		PercentYes float64 `json:"percentYes"`
		PercentNo  float64 `json:"percentNo"`
	}

	if err := initializers.DB.Find(&responses).Error; err != nil {
		log.Println(err)
	}

	Summary.Total = len(responses)
	for _, response := range responses {
		if response.Response && response.Age > 17 {
			Summary.Yes++
		} else if !response.Response && response.Age > 17 {
			Summary.No++
		} else {
			Summary.NA++
		}
	}
	Summary.TotalAdult = Summary.Total - Summary.NA

	// Calculate percentages of valid 'yes/no' answers, rounding to 2 decimal places
	Summary.PercentYes = float64(Summary.Yes) / float64(Summary.TotalAdult) * 100
	Summary.PercentYes = float64(int(Summary.PercentYes*100)) / 100
	Summary.PercentNo = float64(Summary.No) / float64(Summary.TotalAdult) * 100
	Summary.PercentNo = float64(int(Summary.PercentNo*100)) / 100

	c.JSON(http.StatusOK, Summary)
	// free memory for next request
	Summary = struct {
		Total      int     `json:"total"`
		TotalAdult int     `json:"totalAdult"`
		Yes        int     `json:"yes"`
		No         int     `json:"no"`
		NA         int     `json:"na"`
		PercentYes float64 `json:"percentYes"`
		PercentNo  float64 `json:"percentNo"`
	}{}
}

func InterperetData(c *gin.Context) {
	var DataInterperetation struct {
		TotalResponses           int     `json:"totalResponses"`
		TotalAdultResponses      int     `json:"totalAdultResponses"`
		YesResponses             int     `json:"yesResponses"`
		NoResponses              int     `json:"noResponses"`
		NAResponses              int     `json:"naResponses"`
		Ho                       string  `json:"ho"`
		Ha                       string  `json:"ha"`
		Po                       float64 `json:"po"`
		PHat                     float64 `json:"pHat"`
		StdError                 float64 `json:"stdError"`
		ZScore                   float64 `json:"zScore"`
		PValue                   float64 `json:"pvalue"`
		Conclusion               string  `json:"conclusion"`
		ConfidenceInterval       string  `json:"confidenceInterval"`
		ConfidenceIntervalString string  `json:"confidenceIntervalString"`
	}

	const p = .348
	DataInterperetation.Ho = fmt.Sprintf("p = %f", p)
	DataInterperetation.Ha = fmt.Sprintf("p != %f", p)
	DataInterperetation.Po = p
	const alpha = .05
	var responses []models.UserResponse

	if err := initializers.DB.Find(&responses).Error; err != nil {
		log.Println(err)
	}

	DataInterperetation.TotalResponses = len(responses)

	for _, response := range responses {
		if response.Response && response.Age > 17 {
			DataInterperetation.YesResponses++
		} else if !response.Response && response.Age > 17 {
			DataInterperetation.NoResponses++
		} else {
			DataInterperetation.NAResponses++
		}
	}
	DataInterperetation.TotalAdultResponses = DataInterperetation.TotalResponses - DataInterperetation.NAResponses

	DataInterperetation.PHat = float64(DataInterperetation.YesResponses) / float64(DataInterperetation.TotalAdultResponses)
	DataInterperetation.StdError = math.Sqrt(p * (1 - p) / float64(DataInterperetation.TotalAdultResponses))
	DataInterperetation.ZScore = (DataInterperetation.PHat - p) / DataInterperetation.StdError
	DataInterperetation.PValue = stats.NormCdf(DataInterperetation.ZScore, 0, 1)

	if DataInterperetation.PValue < alpha {
		DataInterperetation.Conclusion = "We reject the null hypothesis that the proportion of U.S. adults who sleep less than 7+ hours is 34.8%. We have sufficient evidence to say that the proportion of U.S. adults who sleep less than 7 hours is different than 34.8%."
	} else {
		DataInterperetation.Conclusion = "We fail to reject the null hypothesis that the proportion of U.S. adults who sleep less than 7+ hours is 34.8%. We have insufficient evidence to say that the proportion of U.S. adults who sleep less than 7 hours is different than 34.8%."
	}

	CILeft := DataInterperetation.PHat - 1.96*DataInterperetation.StdError
	CILeft = roundFloat(CILeft, 4)
	CIRight := DataInterperetation.PHat + 1.96*DataInterperetation.StdError
	CIRight = roundFloat(CIRight, 4)
	DataInterperetation.ConfidenceInterval = "[" + fmt.Sprint(CILeft) + ", " + fmt.Sprint(CIRight) + "]"
	DataInterperetation.ConfidenceIntervalString = "We are 95% confident the proportion of U.S. adults who sleep less than 7+ hours is between " + fmt.Sprint(CILeft) + " and " + fmt.Sprint(CIRight) + "."

	c.JSON(http.StatusOK, DataInterperetation)
}

func roundFloat(num float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(num*ratio) / ratio
}
