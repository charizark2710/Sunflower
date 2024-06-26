package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"
)

var PostHistory = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostHistory Start")
	historyBody := model.History{}
	if err := json.Unmarshal(c.Body, &historyBody); err == nil {
		historyObj := model.SysHistory{}
		historyBody.ConvertToDB(&historyObj)
		err := handler.NewHistoryHandler(c.Ctx, &historyObj).Create()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
		}

		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}
}

var GetDetailHistory = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetDetailHistory Start")
	id := c.Param("id")
	historyBody := model.SysHistory{}
	err := handler.NewHistoryHandler(c.Ctx, nil).GetById(id, &historyBody)
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}
	resData := model.History{}
	historyBody.ConvertToJson(&resData)
	utils.Log(LogConstant.Info, "GetDetailHistory End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: historyBody}, nil
}

var UpdateHistory = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "UpdateHistory Start")
	deviceId := c.Param("deviceId")
	rel := model.SysDeviceRel{}
	err := handler.NewDeviceRelHandler(c.Ctx, nil).GetById(deviceId, &rel)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}

	historyBody := model.History{}
	if err := json.Unmarshal(c.Body, &historyBody); err == nil {
		historyBody.Id = rel.HistoryID
		historyModel := model.SysHistory{}
		historyBody.ConvertToDB(&historyModel)
		err := handler.NewHistoryHandler(c.Ctx, &historyModel).Update()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
		}
		utils.Log(LogConstant.Info, "UpdateHistory End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}
}
