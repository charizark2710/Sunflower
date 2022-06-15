package services

import (
	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/gin-gonic/gin"
)

var GetAllDevices = func(c *gin.Context) (model.ResponseTemplate, error) {
	db := model.DbHelper.GetDb()
	utils.Log(LogConstant.Info, "GetAllDevices")
	deviceModel := model.SysDevice{}
	db.Find(&deviceModel)
	return model.ResponseTemplate{HttpCode: 200, Data: model.Devices(deviceModel)}, nil
}

var PostDevice = func(c *gin.Context) (model.ResponseTemplate, error) {
	deviceBody := model.SysDevice{}
	if err := c.BindJSON(&deviceBody); err == nil {
		deviceModel := model.SysDevice(deviceBody)
		err := deviceModel.CreateDevices()
		if err != nil {
			utils.Log(LogConstant.Info, err)
			return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return model.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}
