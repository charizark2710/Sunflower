package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"

	"github.com/gin-gonic/gin"
)

func HistoryRouter(router *gin.Engine) {
	// router.POST(urlconst.PostHistory, controller.Controller)
	router.GET(urlconst.GetDetailHistory, controller.Controller)
	router.PUT(urlconst.PutDetailHistory, controller.Controller)
}
