package controller

import (
	"github.com/charizark2710/Sunflower/RDIPs-BE/constant"
	"github.com/charizark2710/Sunflower/RDIPs-BE/constant/ServiceConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/gin-gonic/gin"
)

func Controller(c *gin.Context) {
	utils.Log(constant.Info, c.FullPath())
	// Get Services
	v := ServiceConst.ServicesMap[c.FullPath()]
	// Handle service
	if fn, ok := v.(func(c *gin.Context) (model.ResponseTemplate, error)); !ok {
		utils.Log(constant.Error, "Wrong type of services functions")
		c.JSON(500, "Wrong type of services functions")
	} else {
		result, err := fn(c)
		if err != nil {
			c.Error(err)
		}
		c.JSON(result.HttpCode, result)
	}
}
