package main

import (
	"os"

	"RDIPs-BE/config"
	LogConstant "RDIPs-BE/constant/LogConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/routers"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.PrepareLog()
	r := gin.Default()

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
	commonModel.Helper.SetDb(db)
	commonModel.Helper.SetMongoDB(mongoDB)
	routers.InitRouter(r)
	r.Run(":" + os.Getenv("API_PORT"))
}
