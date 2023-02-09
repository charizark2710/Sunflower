package main

import (
	"os"

	"github.com/charizark2710/Sunflower/RDIPs-BE/config"
	middleware "github.com/charizark2710/Sunflower/RDIPs-BE/middleware"
	commonModel "github.com/charizark2710/Sunflower/RDIPs-BE/model/common"
	"github.com/charizark2710/Sunflower/RDIPs-BE/routers"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.PrepareLog()
	r := gin.Default()
	r.Use(middleware.SetHeader())
	r.Use(middleware.Validation())

	db, err := config.DbConfig()
	if err == nil {
		commonModel.DbHelper.SetDb(db)
		routers.InitRouter(r)
		r.Run(":" + os.Getenv("PORT"))
	}
}
