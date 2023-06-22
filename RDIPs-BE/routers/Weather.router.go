package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func WeatherRouter(router *gin.Engine) {

	router.Use(middleware.ValidationAPIWeatherKey())
	router.GET(urlconst.GetWeatherForecast, controller.Controller)
}
