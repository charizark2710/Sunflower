package main

import (
	"os"

	"RDIPs-BE/config"
	LogConstant "RDIPs-BE/constant/LogConst"
	middleware "RDIPs-BE/middleware"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/routers"
	"RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.PrepareLog()
	r := gin.Default()
	r.Use(middleware.SetHeader())
	r.Use(middleware.Validation())

	db, err := config.DbConfig()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
		panic(err)
	}
	conn, ch, err := config.RabbitMqConfig()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
		panic(err)
	}
	commonModel.Helper.SetDb(db)
	commonModel.Helper.SetAMQP(conn)
	commonModel.Helper.SetAMQPChannel(ch)
	defer func() {
		conn.Close()
		ch.Close()
	}()
	routers.InitRouter(r)
	go routers.InitAmqpRoutes()
	r.Run(":" + os.Getenv("PORT"))
}
