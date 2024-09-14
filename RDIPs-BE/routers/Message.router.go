package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"

	"github.com/gin-gonic/gin"
)

func MessageRouter(router *gin.Engine) {
	router.POST(urlconst.TransmitMessage, controller.Controller)
}
