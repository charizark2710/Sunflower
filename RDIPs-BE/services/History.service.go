package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

var PostHistory = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostHistory Start")
	historyBody := model.History{}
	if err := c.BindJSON(&historyBody); err == nil {
		historyObj := model.SysHistory{}
		historyBody.ConvertToDB(&historyObj)
		err := handler.Create(&historyObj)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}

var GetDetailHistory = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetDetailHistory Start")
	id := c.Param("id")
	historyBody := model.SysHistory{}
	db := commonModel.DbHelper.GetDb()
	err := db.Where("id = ?", id).First(&historyBody).Error
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	resData := model.History{}
	historyBody.ConvertToJson(&resData)
	utils.Log(LogConstant.Info, "GetDetailHistory End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: historyBody}, nil
}

var UpdateHistory = func(c *gin.Context) (commonModel.ResponseTemplate, error) {
	id := c.Param("id")
	historyBody := model.History{}
	if err := c.BindJSON(&historyBody); err == nil {
		historyModel := model.SysHistory{Id: id}
		err := handler.Update(&historyModel, historyBody)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
}
