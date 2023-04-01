package services

import (
	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/handler"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	commonModel "github.com/charizark2710/Sunflower/RDIPs-BE/model/common"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/gin-gonic/gin"
)

var GetAllPerformances = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetAllPerformances Start")
	var performanceModel []model.SysPerformance
	db := commonModel.DbHelper.GetDb()
	err := db.Find(&performanceModel).Error
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

var PostPerformance = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostPerformance Start")
	performanceBody := model.Performance{}
	if err := c.BindJSON(&performanceBody); err == nil {
		performanceObj := model.SysPerformance{}
		performanceBody.ConvertToDB(&performanceObj)
		err := handler.Create(&performanceObj)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var GetDetailPerformance = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetDetailPerformance Start")
	id := c.Param("id")
	performanceBody := model.SysPerformance{}
	db := commonModel.DbHelper.GetDb()
	err := db.Where("id = ?", id).First(&performanceBody).Error
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	resData := model.Performance{}
	performanceBody.ConvertToJson(&resData)
	utils.Log(LogConstant.Info, "GetDetailPerformance End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: performanceBody}, nil
}

var PutPerformance = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "UpdatePerformance Start")
	id := c.Param("id")
	performanceBody := model.Performance{}
	if err := c.BindJSON(&performanceBody); err == nil {
		performanceModel := model.SysPerformance{Id: id}
		err := handler.Update(&performanceModel, performanceBody)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}
