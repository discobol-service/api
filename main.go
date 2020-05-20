package main

import "github.com/gin-gonic/gin"


func main() {
	router := gin.Default()

	router.GET("/v1/ping", GETPingV1)
	router.GET("/v1/feed", GETFeedV1)
	router.GET("/v1/click", GETClickV1)

	router.Run()
}


func GETPingV1(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GETFeedV1(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "feed!",
	})
}

func GETClickV1(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "click",
	})
}