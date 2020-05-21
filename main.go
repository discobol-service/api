package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type discount struct {
	Title       string
	Description string
	Thumb       string
	Type        int64
	URL         string
	AmountFrom  float64
	AmountTo    float64
	Currency    string
}

type recsRequest struct {
	RegionID int64
	UserID   int64
}

func main() {
	router := gin.Default()

	router.GET("/v1/recs", recsV1)
	router.POST("/v1/click", clickV1)

	router.Run()
}

// Взять у бандита рекомендации,
// соорудить из этого ответ в жсоне:
func recsV1(c *gin.Context) {

	discounts, err := getLatestDiscounts(100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	if len(discounts) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "no discounts",
		})
		return
	}

	resp, err := http.Get("http://localhost:4000/recs/100")
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(404, gin.H{
			"message": "recs not found",
		})
		return
	}

	for true {
		bs := make([]byte, 1014)
		n, err := resp.Body.Read(bs)
		log.Println("message from bandit: " + string(bs[:n]))

		if n == 0 || err != nil {
			break
		}
	}

	c.JSON(200, gin.H{
		"message": "feed!",
	})
}

func clickV1(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "click",
	})
}

func getLatestDiscounts(limit int) ([]*discount, error) {
	list := make([]*discount, limit)

	return list, nil
}
