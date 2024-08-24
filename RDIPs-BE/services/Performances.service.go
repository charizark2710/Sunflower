package services

import (
	"RDIPs-BE/constant"
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var GetAllPerformances = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetAllPerformances Start")
	var performanceModel []model.SysPerformance
	err := handler.NewPerformanceHandler(c.Ctx, nil).Read(&performanceModel)
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}
	resData := make([]model.Performance, len(performanceModel))
	for i, performance := range performanceModel {
		performance.ConvertToJson(&resData[i])
	}
	utils.Log(LogConstant.Info, "GetAllPerformances End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: resData}, nil
}

var PostPerformance = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostPerformance Start")
	performanceBody := model.Performance{}
	if err := json.Unmarshal(c.Body, &performanceBody); err == nil {
		performanceObj := model.SysPerformance{}
		performanceBody.ConvertToDB(&performanceObj)
		err := handler.NewPerformanceHandler(c.Ctx, &performanceObj).Create()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}
}

var GetDetailPerformance = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetDetailPerformance Start")
	id := c.Param("id")
	fromQ := c.Query("from")
	toQ := c.Query("to")
	filterByStr := c.Query("filterDocBy")
	var opts []map[string]interface{}
	if filterByStr != "" {
		if err := json.Unmarshal([]byte(filterByStr), &opts); err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
		}
	}

	var from time.Time
	var to time.Time
	var err error
	if fromQ != "" {
		from, err = time.Parse(time.RFC3339, c.Query("from"))
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 400, Data: nil, Message: err.Error()}, err
		}
	}
	if toQ != "" {
		to, err = time.Parse(time.RFC3339, c.Query("to"))
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 400, Data: nil, Message: err.Error()}, err
		}
	}

	var filterOpt map[string]interface{}
	if toQ != "" || fromQ != "" {
		filterOpt = map[string]interface{}{
			"filterBytime": map[string]time.Time{
				"from": from,
				"to":   to,
			},
		}
		opts = append(opts, filterOpt)
	}

	docLimit, err := strconv.Atoi(c.Query("docLimit"))
	if err != nil {
		docLimit = 100
	}

	limitOpt := map[string]interface{}{
		"limit": docLimit,
	}

	opts = append(opts, limitOpt)

	performanceBody := model.SysPerformance{}
	err = handler.NewPerformanceHandler(c.Ctx, nil).GetById(id, &performanceBody, opts...)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}
	resData := model.Performance{}
	performanceBody.ConvertToJson(&resData)
	utils.Log(LogConstant.Info, "GetDetailPerformance End")
	return commonModel.ResponseTemplate{HttpCode: 200, Data: performanceBody}, nil
}

var PutPerformance = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "UpdatePerformance Start")
	deviceId := c.Param("deviceId")
	rel := model.SysDeviceRel{}
	err := handler.NewDeviceRelHandler(c.Ctx, nil).GetById(deviceId, &rel)
	if err == gorm.ErrRecordNotFound {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 404, Data: nil, Message: err.Error()}, err
	} else if err != nil {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}

	performanceBody := model.Performance{}
	if err := json.Unmarshal(c.Body, &performanceBody); err == nil {
		if c.Header[constant.REQUEST_TYPE_HEADER] == nil &&
			performanceBody.Payload != nil {
			errMsg := "cannot update performance if request is not from mqtt"
			return commonModel.ResponseTemplate{HttpCode: 400, Data: nil, Message: errMsg}, fmt.Errorf(errMsg)
		}
		performanceBody.Id = rel.PerformanceID
		performanceModel := model.SysPerformance{}
		performanceModel.UpdatedAt = time.Now()
		performanceBody.ConvertToDB(&performanceModel)
		err := handler.NewPerformanceHandler(c.Ctx, &performanceModel).Update()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
		}
		utils.Log(LogConstant.Info, "UpdatePerformance End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
	}
}
