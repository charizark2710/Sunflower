package controller

import (
	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/constant/ServiceConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/gin-gonic/gin"
)

func Controller(c *gin.Context) {
	// Get Services
	v := ServiceConst.ServicesMap[c.Request.Method+c.FullPath()]
	// Handle service
	if fn, ok := v.(func(c *gin.Context) (model.ResponseTemplate, error)); !ok {
		utils.Log(LogConstant.Error, "Wrong type of services functions")
		c.JSON(500, "Wrong type of services functions")
	} else {
		result, err := fn(c)
		result.SetMessage(err.Error())
		result.SetError(err)
		if err != nil {
			c.Error(err)
		}
		c.JSON(result.HttpCode, result)
	}
}
