package middleware

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ValidatorConst"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

func Validator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Services
		v := ValidatorConst.ValidatorsMap[c.Request.Method+c.FullPath()]
		// Handle validate
		if fn, ok := v.(func(c *gin.Context) error); !ok {
			utils.Log(LogConstant.Error, "Wrong type of validator functions")
			c.JSON(500, "Wrong type of validator functions")
		} else {
			err := fn(c)
			if err != nil {
				utils.Log(LogConstant.Error, err)
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
