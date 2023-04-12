package controller

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

func Controller(c *gin.Context) {
	// Get Services
	v := ServiceConst.ServicesMap[c.Request.Method+c.FullPath()]
	// Handle service
	if fn, ok := v.(func(c *gin.Context) (commonModel.ResponseTemplate, error)); !ok {
		utils.Log(LogConstant.Error, "Wrong type of services functions")
		c.JSON(500, "Wrong type of services functions")
	} else {
		result, err := fn(c)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			result.SetMessage(err.Error())
		}
		c.JSON(result.HttpCode, result)
	}
}
