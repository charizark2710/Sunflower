package controller

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"context"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

func Controller(c *gin.Context) {
	// Prepare Services
	fn := ServiceConst.ServicesMap[c.Request.Method+c.FullPath()]
	bodyAsByteArray, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		c.JSON(500, err)
	}
	serviceCtx := commonModel.ServiceContext{Ctx: context.Background(), Mu: sync.Mutex{}, ServiceModel: commonModel.ServiceModel{
		Body: bodyAsByteArray,
	}}
	setQueryAndParam(c, &serviceCtx)

	// Handle service
	result, err := fn(&serviceCtx)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		result.Error = err
		result.SetMessage(err.Error())
	}
	c.JSON(result.HttpCode, result)
}

func setQueryAndParam(c *gin.Context, service *commonModel.ServiceContext) {
	service.InitParamsAndQueries()

	queries := c.Request.URL.Query()
	params := c.Params

	for qk, qv := range queries {
		service.SetQuery(qk, qv[0])
	}

	for _, p := range params {
		service.SetParam(p.Key, p.Value)
	}
}
