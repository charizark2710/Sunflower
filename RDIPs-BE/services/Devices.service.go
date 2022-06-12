package services

import (
	"github.com/charizark2710/Sunflower/RDIPs-BE/constant"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/gin-gonic/gin"
)

var GetAllDevices = func(c *gin.Context) (model.ResponseTemplate, error) {
	utils.Log(constant.Info, "Test")
	data := model.Device{Id: "123", Name: "Test"}
	return model.ResponseTemplate{HttpCode: 200, Data: data}, nil
}
