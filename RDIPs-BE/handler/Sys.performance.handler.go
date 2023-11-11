package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/model"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

type PerformanceHandler interface {
	CommonHandler
}

type performanceHandler struct {
	performanceBody *model.SysPerformance
	*commonHandler
}

func NewPerformanceHandler(c *gin.Context, performanceModel *model.SysPerformance) PerformanceHandler {
	commonHanlerInstance := newCommonHandler(c)
	commonStruct := commonHanlerInstance.(*commonHandler)
	return &performanceHandler{commonHandler: commonStruct, performanceBody: performanceModel}
}

func (p *performanceHandler) Read(performanceRes interface{}) error {
	return p.db.Find(performanceRes).Error
}

func (p *performanceHandler) Create() error {
	performanceObj := p.performanceBody
	return p.db.Save(&performanceObj).Error
}

func (p *performanceHandler) GetById(id string, performanceResponse interface{}) error {
	return p.db.Where("id = ?", id).First(performanceResponse).Error
}

func (p *performanceHandler) Update() error {
	err := p.GetById(p.performanceBody.Id, &model.SysPerformance{})
	if err != nil {
		utils.Log(LogConstant.Error, "Cannot find performance with Id = "+p.performanceBody.Id, err)
		return err
	}
	return p.db.Updates(p.performanceBody).Error
}
