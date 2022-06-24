package main

import (
	"os"

	"github.com/charizark2710/Sunflower/RDIPs-BE/config"
	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	middleware "github.com/charizark2710/Sunflower/RDIPs-BE/middleware"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	"github.com/charizark2710/Sunflower/RDIPs-BE/routers"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	utils.PrepareLog()
	if err != nil {
		utils.Log(LogConstant.Error, "Error loading .env file")
	}

	r := gin.Default()
	r.Use(middleware.SetHeader())
	r.Use(middleware.Validation())

	db, err := config.DbConfig()
	if err == nil {
		model.DbHelper.SetDb(db)
		routers.InitRouter(r)
		r.Run(":" + os.Getenv("PORT"))
	}
}
