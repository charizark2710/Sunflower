package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"

	"gorm.io/gorm"
)

var GetAllDevices = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetAllDevices Start")
	var deviceModel []model.SysDevices
	db := commonModel.Helper.GetDb()
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

var PostDevice = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostDevice Start")
	defer utils.Log(LogConstant.Info, "PostDevice End")
	deviceBody := model.Devices{}
	if err := json.Unmarshal(c.Body, &deviceBody); err == nil {
		deviceObj := model.SysDevices{}
		deviceBody.ConvertToDB(&deviceObj)

		db := commonModel.Helper.GetDb()
		err := db.Transaction(func(tx *gorm.DB) error {
			utils.Log(LogConstant.Info, "Create Device Start")
			if err = handler.CreateWithTx(&deviceObj, tx); err != nil {
				return err
			}

			historyObj := model.SysHistory{
				LogPath: deviceObj.Name + "/",
			}
			utils.Log(LogConstant.Info, "Create History Start")
			if err := handler.CreateWithTx(&historyObj, tx); err != nil {
				return err
			}

			performanceObj := model.SysPerformance{
				DocumentName: deviceObj.Name,
			}
			utils.Log(LogConstant.Info, "Create Performance Start")
			if err := handler.CreateWithTx(&performanceObj, tx); err != nil {
				return err
			}

			deviceRelObj := model.SysDeviceRel{
				DeviceID:      deviceObj.Id,
				PerformanceID: performanceObj.Id,
				HistoryID:     historyObj.Id,
			}
			utils.Log(LogConstant.Info, "Create device_rel Start")
			if err := handler.CreateWithTx(&deviceRelObj, tx); err != nil {
				return err
			}
			return nil
		})

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
	db := commonModel.Helper.GetDb()
	var err error
	if detail == "true" {
		err = db.Where("id = ? AND status != ?", id, model.Disable).Preload("DeviceRel").Preload("DeviceRel.History").
			Preload("DeviceRel.Performance").First(&deviceBody).Error
	} else {
		err = db.Where("id = ? AND status != ?", id, model.Disable).Preload("DeviceRel").First(&deviceBody).Error
	}

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

var DeleteDevice = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "DeleteDevice Start")
	defer utils.Log(LogConstant.Info, "DeleteDevice End")
	id := c.Param("id")
	err := handler.Update(&model.SysDevices{Id: id}, model.SysDevices{Status: model.Disable})
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil, Message: ""}, err
}
