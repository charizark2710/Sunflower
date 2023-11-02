package handler

import (
	"RDIPs-BE/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeviceRelHandler interface {
	CommonHandler
	GetDeviceRelByDeviceId() error
}

type deviceRelHandler struct {
	deviceRelBody *model.SysDeviceRel
	*commonHandler
}

func NewDeviceRelHandler(c *gin.Context, deviceRelModel *model.SysDeviceRel) DeviceRelHandler {
	commonHanlerInstance := newCommonHandler(c)
	commonStruct := commonHanlerInstance.(*commonHandler)
	return &deviceRelHandler{commonHandler: commonStruct, deviceRelBody: deviceRelModel}
}

func (dRel *deviceRelHandler) GetDeviceRelByDeviceId() error {
	return dRel.db.Where("device_id = ?", dRel.deviceRelBody.DeviceID).First(dRel.deviceRelBody).Preload("sys_history", func(db *gorm.DB) *gorm.DB {
		return db.Order("sunflower.sys_history.log_path Desc").Limit(1)
	}).Error
}
