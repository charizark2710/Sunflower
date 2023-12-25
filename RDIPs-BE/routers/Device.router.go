package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"

	"github.com/gin-gonic/gin"
)

func DevicesRouter(router *gin.Engine) {
	router.GET(urlconst.GetAllDevices, controller.Controller)
	router.POST(urlconst.PostDevice, controller.Controller)
	router.GET(urlconst.GetDetailDevice, controller.Controller)
	router.PUT(urlconst.PutDetailDevice, controller.Controller)
	router.DELETE(urlconst.DeleteDevice, controller.Controller)
	router.GET(urlconst.GetLogOfDevice, controller.Controller)
	router.POST(urlconst.PostLogOfDevice, controller.Controller)
}
