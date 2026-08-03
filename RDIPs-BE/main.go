package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"RDIPs-BE/config"
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/routers"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Extract request details after handlers run
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		logMessage := fmt.Sprintf("| %3d | %13v | %15s | %-7s %s %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)

		// Route to your logger based on HTTP status code
		switch {
		case statusCode >= 500:
			utils.Log(LogConstant.Error, logMessage)
		case statusCode >= 400:
			utils.Log(LogConstant.Warning, logMessage)
		default:
			utils.Log(LogConstant.Info, logMessage)
		}
	}
}

func main() {
	utils.PrepareLog()
	r := gin.New()
	r.Use(GinLogger())
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		} else {
			c.String(http.StatusInternalServerError, "Internal server error")
		}
	}))

	err := handler.ReadRootCACert()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	db, err := config.DbConfig()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	mongoDB, err := config.MongoConfig()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	err = config.RabbitMqConfig()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	err = config.KeycloakConfig()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	commonModel.Factory.SetDB(commonModel.Postgres, db)
	commonModel.Factory.SetDB(commonModel.MongoDB, mongoDB)

	routers.InitRouter(r)
	r.Run(":" + os.Getenv("API_PORT"))
}
