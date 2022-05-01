package main

import (
	"log"
	"os"

	"github.com/charizark2710/Sunflower/RDIPs-BE/constant"
	middleware "github.com/charizark2710/Sunflower/RDIPs-BE/middleware"
	"github.com/charizark2710/Sunflower/RDIPs-BE/routers"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	utils.PrepareLog()
	utils.Log(constant.Info, "start sunflower server")
	r := gin.Default()
	r.Use(middleware.SetHeader())
	r.Use(middleware.Validation())
	routers.InitRouter(r)
	r.Run(":" + os.Getenv("PORT"))
}
