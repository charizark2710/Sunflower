package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/model"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

type HistoryHandler interface {
	CommonHandler
}

type historyHandler struct {
	historyBody *model.SysHistory
	*commonHandler
}

func NewHistoryHandler(c *gin.Context, historyModel *model.SysHistory) HistoryHandler {
	commonHanlerInstance := newCommonHandler(c)
	commonStruct := commonHanlerInstance.(*commonHandler)
	return &historyHandler{commonHandler: commonStruct, historyBody: historyModel}
}

func (h *historyHandler) Create() error {
	historyObj := h.historyBody
	return h.db.Save(&historyObj).Error
}

func (h *historyHandler) GetById(id string, historyResponse interface{}) error {
	return h.db.Where("id = ?", id).First(historyResponse).Error
}

func (h *historyHandler) Update() error {
	err := h.GetById(h.historyBody.Id, &model.SysHistory{})
	if err != nil {
		utils.Log(LogConstant.Error, "Cannot find history with Id = "+h.historyBody.Id, err)
		return err
	}
	return h.db.Updates(h.historyBody).Error
}
