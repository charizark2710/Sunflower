package routers

import (
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Use(middleware.ValidatorMiddleware())
	router.Use(middleware.SetFilter())
	DevicesRouter(router)
	PerformanceRouter(router)
	HistoryRouter(router)
	WeatherRouter(router)
	UserRouter(router)
	GroupRouter(router)
}
