package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

var GetAllDevices = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetAllDevices Start")
	var deviceModel []model.SysDevices
	db := commonModel.DbHelper.GetDb()
	err := db.Where("status != ?", model.Disable).Preload("Parent").Find(&deviceModel).Error
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	resData := make([]model.Devices, len(deviceModel))
	for i, device := range deviceModel {
		device.ConvertToJson(&resData[i])
	}
	utils.Log(LogConstant.Info, "GetAllDevices End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: resData}, nil
}

var PostDevice = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	deviceBody := model.Devices{}
	// db := commonModel.DbHelper.GetDb()
	if err := c.BindJSON(&deviceBody); err == nil {
		deviceObj := model.SysDevices{}
		deviceBody.ConvertToDB(&deviceObj)
		err := handler.Create(&deviceObj)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var GetDetailDevice = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	id := c.Param("id")
	deviceBody := model.SysDevices{}
	db := commonModel.DbHelper.GetDb()
	err := db.Where("id = ? AND status != ?", id, model.Disable).First(&deviceBody).Error
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	resData := model.Devices{}
	deviceBody.ConvertToJson(&resData)
	return commonModel.ResponseTemplate{HttpCode: 200, Data: deviceBody}, nil
}

var UpdateDevice = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	id := c.Param("id")
	deviceBody := model.Devices{}
	if err := c.BindJSON(&deviceBody); err == nil {
		deviceModel := model.SysDevices{Id: id}
		err := handler.Update(&deviceModel, deviceBody)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var DeleteDevice = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	id := c.Param("id")
	err := handler.Update(&model.SysDevices{Id: id}, model.SysDevices{Status: model.Disable})
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil, Message: ""}, err
}
