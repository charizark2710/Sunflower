package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func GroupRouter(router *gin.Engine) {
	authRouter := router.Group("")
	{
		authRouter.Use(middleware.CheckClientTokenValidation())
		authRouter.GET(urlconst.GetKeycloakGroups, controller.Controller)
		authRouter.GET(urlconst.GetKeycloakGroupById, controller.Controller)
		authRouter.DELETE(urlconst.DeleteKeycloakGroup, controller.Controller)
		authRouter.POST(urlconst.PostKeycloakGroup, controller.Controller)
		authRouter.PUT(urlconst.PutKeycloakGroup, controller.Controller)
	}

}
