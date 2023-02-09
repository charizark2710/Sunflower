package controller

import (
	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/constant/ServiceConst"
	commonModel "github.com/charizark2710/Sunflower/RDIPs-BE/model/common"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
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
			result.SetMessage(err.Error())
			result.SetError(err)
			c.Error(err)
			return
		}
		c.JSON(result.HttpCode, result)
	}
}
