package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"

	"github.com/gin-gonic/gin"
)

func WeatherRouter(router *gin.Engine) {
	router.GET(urlconst.GetWeatherNext14Days, controller.Controller)
}
