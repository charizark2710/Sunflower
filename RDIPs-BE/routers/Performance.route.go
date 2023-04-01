package routers

import (
	urlconst "github.com/charizark2710/Sunflower/RDIPs-BE/constant/URLConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/controller"
	"github.com/gin-gonic/gin"
)

func PerformanceRouter(router *gin.Engine) {
	router.GET(urlconst.GetAllPerformances, controller.Controller)
	router.POST(urlconst.PostPerformance, controller.Controller)
	router.GET(urlconst.GetDetailPerformance, controller.Controller)
	router.PUT(urlconst.PutDetailPerformance, controller.Controller)
}
