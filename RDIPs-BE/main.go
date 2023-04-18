package main

import (
	"os"

	"RDIPs-BE/config"
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
	r.Use(middleware.Validator())

	db, err := config.DbConfig()
	if err == nil {
		commonModel.DbHelper.SetDb(db)
		routers.InitRouter(r)
		r.Run(":" + os.Getenv("PORT"))
	}
}
