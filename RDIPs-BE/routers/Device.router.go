package routers

import (
	urlconst "github.com/charizark2710/Sunflower/RDIPs-BE/constant/URLConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/controller"
	"github.com/gin-gonic/gin"
)

func DevicesRouter(router *gin.Engine) {
	router.GET(urlconst.GetAllDevices, controller.Controller)
	router.POST(urlconst.PostDevice, controller.Controller)
	router.GET(urlconst.GetDetailDevice, controller.Controller)
	router.PUT(urlconst.PutDetailDevice, controller.Controller)
	router.DELETE(urlconst.DeleteDevice, controller.Controller)

	router.POST(urlconst.PostHistory, controller.Controller)
	router.GET(urlconst.GetDetailHistory, controller.Controller)
	router.GET(urlconst.GetHistoriesOfDevice, controller.Controller)
}
