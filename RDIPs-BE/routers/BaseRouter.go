package routers

import (
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	// For health check only
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})
	router.Use(middleware.SetHeader())
	router.Use(middleware.Validation())
	router.Use(middleware.ValidatorMiddleware())
	router.Use(middleware.SetFilter())
	DevicesRouter(router)
	PerformanceRouter(router)
	HistoryRouter(router)
	WeatherRouter(router)
	KeycloakRouter(router)
}
