package routers

import (
	"github.com/gin-gonic/gin"
)

func BaseRouter(router *gin.Engine) {
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
