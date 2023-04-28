package middleware

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ValidatorConst"
	"RDIPs-BE/utils"
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Services
		v := ValidatorConst.ValidatorsMap[c.Request.Method+c.FullPath()]
		// Handle validate
		if v == nil {
			utils.Log(LogConstant.Info, "This API do not have validator")
		} else {
			jsonBody, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}
			switch c.Request.Method {
			case "POST":
				err = v.Post(jsonBody)
			case "PUT":
				err = v.Put(jsonBody)
			}

			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
				return
			}
			c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
		}

		c.Next()
	}
}
