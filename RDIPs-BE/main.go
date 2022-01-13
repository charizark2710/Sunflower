package main

import (
	"log"
	"os"

	mimddlewere "github.com/charizark2710/Automate-Garden/RDIPs-BE/middlewere"
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
	r.Use(mimddlewere.SetHeader())
	r.Use(mimddlewere.Validation())
	routers.BaseRouter(r)
	r.Run(":" + os.Getenv("PORT"))
}
