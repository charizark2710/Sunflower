package services

import (
	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/handler"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/gin-gonic/gin"
)

var GetAllDevices = func(c *gin.Context) (model.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetAllDevices")
	var deviceModel []model.SysDevice
	err := handler.Read(&deviceModel)
	if err != nil {
		return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	var resData []model.Devices
	for _, device := range deviceModel {
		resData = append(resData, model.Devices{Id: device.Id, Name: device.Name})
	}
	return model.ResponseTemplate{HttpCode: 200, Data: resData}, nil
}

var PostDevice = func(c *gin.Context) (model.ResponseTemplate, error) {
	deviceBody := model.Devices{}
	if err := c.BindJSON(&deviceBody); err == nil {
		deviceModel := model.SysDevice{Id: deviceBody.Id, Name: deviceBody.Name}
		err := handler.Create(&deviceModel)
		if err != nil {
			utils.Log(LogConstant.Info, err)
			return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return model.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var GetDetailDevice = func(c *gin.Context) (model.ResponseTemplate, error) {
	id := c.Param("id")
	deviceBody := model.SysDevice{}
	if err := handler.ReadDetail(&deviceBody, id); err == nil {
		return model.ResponseTemplate{HttpCode: 200, Data: model.Devices{Id: deviceBody.Id, Name: deviceBody.Name}}, nil
	} else {
		return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var UpdateDevice = func(c *gin.Context) (model.ResponseTemplate, error) {
	id := c.Param("id")
	deviceBody := model.Devices{}
	if err := c.BindJSON(&deviceBody); err == nil {
		deviceModel := model.SysDevice{Id: id}
		err := handler.Update(&deviceModel, deviceBody)
		if err != nil {
			utils.Log(LogConstant.Info, err)
			return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return model.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var DeleteDevice = func(c *gin.Context) (model.ResponseTemplate, error) {
	id := c.Param("id")
	err := handler.Update(&model.SysDevice{Id: id}, model.SysDevice{Delete: true})
	if err != nil {
		return model.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return model.ResponseTemplate{HttpCode: 200, Data: nil, Message: ""}, err
}
