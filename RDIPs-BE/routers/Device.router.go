package routers

import (
	"github.com/charizark2710/Sunflower/RDIPs-BE/controller"
	"github.com/gin-gonic/gin"
)

func DeviceRouter(router *gin.Engine) {
	router.GET("/devices", controller.GetAllDevicesController)
}
