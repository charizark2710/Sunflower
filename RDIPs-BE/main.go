package main

import (
	mimddlewere "github.com/charizark2710/Automate-Garden/RDIPs-BE/middlewere"
	"github.com/charizark2710/Automate-Garden/RDIPs-BE/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(mimddlewere.SetHeader())
	routers.BaseRouter(r)
	r.Run(":8080")
}
