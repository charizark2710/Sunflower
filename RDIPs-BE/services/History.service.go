package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	id := c.Param("deviceId")
	historyBody := model.SysHistory{}
	db := commonModel.Helper.GetDb()
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
	deviceId := c.Param("deviceId")
	amqp := c.Query("amqp")
	db := commonModel.Helper.GetDb()
	rel := model.DeviceRel{DeviceID: deviceId}
	historyModel := model.SysHistory{}
	db.First(&rel).Preload("sys_history", func(db *gorm.DB) *gorm.DB {
		return db.Order("sunflower.sys_history.log_path Desc").Limit(1)
	})
	historyBody := model.History{}
	if err := c.BindJSON(&historyBody); err == nil {
		if amqp == "true" {
			go func() {
				fileIO := handler.FileIO{Name: rel.History.LogPath}
				fileIO.Write(time.Now(), []byte(historyBody.Payload))
			}()
		}
		historyModel = model.SysHistory{Id: rel.HistoryID}
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
