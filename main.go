package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type Rec struct {
	Title       string
	Description string
	Thumb       string
	Type        int64
	URL         string
	AmountFrom  float64
	AmountTo    float64
	Currency    string
}


func main() {
	router := gin.Default()

	router.GET("/v1/recs", recsV1)
	router.GET("/v1/click", clickV1)

	router.Run()
}

// Взять у бандита рекомендации,
// соорудить из этого ответ в жсоне:
func recsV1(c *gin.Context) {
	resp, err := http.Get("http://localhost:4000/recs/300")
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error",
		})
		return
	}
	defer resp.Body.Close()

	for true {
        bs := make([]byte, 1014)
        n, err := resp.Body.Read(bs)
        log.Println("message from bandit: " + string(bs[:n]))

        if n == 0 || err != nil{
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