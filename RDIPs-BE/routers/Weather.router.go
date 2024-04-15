package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func WeatherRouter(router *gin.Engine) {
	router.GET(urlconst.GetWeatherForecast, middleware.ValidationAPIWeatherKey(), controller.Controller)
}
