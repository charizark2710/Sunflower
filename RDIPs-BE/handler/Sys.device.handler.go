package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/model"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeviceHandler interface {
	CommonHandler
	ReadDetail(isDetail bool, id string) error
}

type deviceHandler struct {
	deviceBody *model.SysDevices
	*commonHandler
}

func NewDeviceHandler(c *gin.Context, deviceModel *model.SysDevices) DeviceHandler {
	commonHanlerInstance := newCommonHandler(c)
	commonStruct := commonHanlerInstance.(*commonHandler)
	return &deviceHandler{commonHandler: commonStruct, deviceBody: deviceModel}
}

func (d *deviceHandler) Read(devicesRes interface{}) error {
	return d.db.Where("status != ?", model.Disable).Preload("Parent").Find(devicesRes).Error
}

func (d *deviceHandler) Create() error {
	deviceObj := d.deviceBody
	return d.db.Transaction(func(tx *gorm.DB) error {
		utils.Log(LogConstant.Info, "Create Device Start")
		if err := tx.Create(&deviceObj).Error; err != nil {
			return err
		}

		historyObj := model.SysHistory{
			LogPath: deviceObj.Name + "/",
		}
		utils.Log(LogConstant.Info, "Create History Start")
		if err := tx.Create(&historyObj).Error; err != nil {
			return err
		}

		performanceObj := model.SysPerformance{
			DocumentName: deviceObj.Name,
		}
		utils.Log(LogConstant.Info, "Create Performance Start")
		if err := tx.Create(&performanceObj).Error; err != nil {
			return err
		}

		deviceRelObj := model.SysDeviceRel{
			DeviceID:      deviceObj.Id,
			PerformanceID: performanceObj.Id,
			HistoryID:     historyObj.Id,
		}
		utils.Log(LogConstant.Info, "Create device_rel Start")
		if err := tx.Create(&deviceRelObj).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *deviceHandler) ReadDetail(isDetail bool, id string) error {
	var err error
	if isDetail {
		err = d.db.Where("id = ? AND status != ?", id, model.Disable).Preload("DeviceRel").Preload("DeviceRel.History").
			Preload("DeviceRel.Performance").First(d.deviceBody).Error
	} else {
		err = d.db.Where("id = ? AND status != ?", id, model.Disable).Preload("DeviceRel").First(d.deviceBody).Error
	}
	return err
}

func (d *deviceHandler) GetById(id string, response interface{}) error {
	return d.db.Where("id = ?", id).First(response).Error
}

func (d *deviceHandler) Update() error {
	err := d.GetById(d.deviceBody.Id, &model.SysDevices{})
	if err != nil {
		utils.Log(LogConstant.Error, "cannot find device with ID = "+d.deviceBody.Id, err)
		return err
	}
	return d.db.Updates(d.deviceBody).Error
}
