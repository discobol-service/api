package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/discobol-service/api/controller"
	"github.com/discobol-service/api/controller/v1"
)

func Routes(e *gin.Engine) {
	e.GET("/ping", controller.Ping)

	e.GET("/v1/discounts", v1.GetDiscounts)
}