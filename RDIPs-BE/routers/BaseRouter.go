package routers

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	DevicesRouter(router)
}
