package services

import (
	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/handler"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	commonModel "github.com/charizark2710/Sunflower/RDIPs-BE/model/common"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
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
