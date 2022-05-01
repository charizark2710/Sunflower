package controller

import (
	"github.com/charizark2710/Sunflower/RDIPs-BE/constant"
	"github.com/charizark2710/Sunflower/RDIPs-BE/services"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/gin-gonic/gin"
)

func GetAllDevicesController(c *gin.Context) {
	utils.Log(constant.Info, "GetAllDevicesController")
	// Handle service
	result, err := services.GetAllDevices()
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, result)
}
