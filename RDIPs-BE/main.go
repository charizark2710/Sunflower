package main

import (
	mimddlewere "./middlewere"
	"./routers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(mimddlewere.SetHeader())
	routers.BaseRouter(r)
	r.Run(":8080")
}
