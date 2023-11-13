package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"
)

var GetAllPerformances = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetAllPerformances Start")
	var performanceModel []model.SysPerformance
	err := handler.NewPerformanceHandler(c.Ctx, nil).Read(&performanceModel)
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	resData := make([]model.Performance, len(performanceModel))
	for i, performance := range performanceModel {
		performance.ConvertToJson(&resData[i])
	}
	utils.Log(LogConstant.Info, "GetAllPerformances End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: resData}, nil
}

var PostPerformance = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostPerformance Start")
	performanceBody := model.Performance{}
	if err := json.Unmarshal(c.Body, &performanceBody); err == nil {
		performanceObj := model.SysPerformance{}
		performanceBody.ConvertToDB(&performanceObj)
		err := handler.NewPerformanceHandler(c.Ctx, &performanceObj).Create()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var GetDetailPerformance = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetDetailPerformance Start")
	id := c.Param("id")
	performanceBody := model.SysPerformance{}
	err := handler.NewPerformanceHandler(c.Ctx, nil).GetById(id, &performanceBody)
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	resData := model.Performance{}
	performanceBody.ConvertToJson(&resData)
	utils.Log(LogConstant.Info, "GetDetailPerformance End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: performanceBody}, nil
}

var PutPerformance = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "UpdatePerformance Start")
	deviceId := c.Param("deviceId")
	rel := model.SysDeviceRel{}
	err := handler.NewDeviceRelHandler(c.Ctx, nil).GetById(deviceId, &rel)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}

	performanceBody := model.Performance{}
	if err := json.Unmarshal(c.Body, &performanceBody); err == nil {
		performanceBody.Id = rel.PerformanceID
		performanceModel := model.SysPerformance{}
		performanceBody.ConvertToDB(&performanceModel)
		err := handler.NewPerformanceHandler(c.Ctx, &performanceModel).Update()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		utils.Log(LogConstant.Info, "UpdatePerformance End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}
