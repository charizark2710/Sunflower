package main

import (
	"log"
	"os"

	middleware "github.com/charizark2710/Automate-Garden/RDIPs-BE/middleware"
	"github.com/charizark2710/Automate-Garden/RDIPs-BE/routers"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	r.Use(middleware.SetHeader())
	r.Use(middleware.Validation())
	routers.BaseRouter(r)
	r.Run(":" + os.Getenv("PORT"))
}
