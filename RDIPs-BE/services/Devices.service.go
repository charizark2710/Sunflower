package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"
)

var GetAllDevices = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetAllDevices Start")
	var deviceModel []model.SysDevices
	err := handler.NewDeviceHandler(c.Ctx, nil).Read(&deviceModel)
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

var PostDevice = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostDevice Start")
	defer utils.Log(LogConstant.Info, "PostDevice End")
	deviceBody := model.Devices{}
	if err := json.Unmarshal(c.Body, &deviceBody); err == nil {
		deviceObj := model.SysDevices{}
		deviceBody.ConvertToDB(&deviceObj)
		err := handler.NewDeviceHandler(c.Ctx, &deviceObj).Create()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: map[string]interface{}{"Id": deviceObj.Id}}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var GetDetailDevice = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetDetailDevice Start")
	defer utils.Log(LogConstant.Info, "GetDetailDevice End")
	id := c.Param("id")
	detail := c.Query("detail")
	deviceBody := model.SysDevices{}
	err := handler.NewDeviceHandler(c.Ctx, &deviceBody).ReadDetail(detail == "true", id)
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	resData := model.Devices{}
	deviceBody.ConvertToJson(&resData)
	return commonModel.ResponseTemplate{HttpCode: 200, Data: resData}, nil
}

var UpdateDevice = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "UpdateDevice Start")
	defer utils.Log(LogConstant.Info, "UpdateDevice End")
	id := c.Param("id")
	deviceBody := model.Devices{}
	if err := json.Unmarshal(c.Body, &deviceBody); err == nil {
		deviceBody.Id = id
		deviceModel := model.SysDevices{}
		deviceBody.ConvertToDB(&deviceModel)
		err := handler.NewDeviceHandler(c.Ctx, &deviceModel).Update()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var DeleteDevice = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "DeleteDevice Start")
	defer utils.Log(LogConstant.Info, "DeleteDevice End")
	id := c.Param("id")
	err := handler.NewDeviceHandler(c.Ctx, &model.SysDevices{Id: id, Status: model.Disable}).Update()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil, Message: ""}, err
}
