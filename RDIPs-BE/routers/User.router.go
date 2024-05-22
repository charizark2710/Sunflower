package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	authRouter := router.Group("")
	{
		router.POST(urlconst.PostLogin, controller.Controller)
		authRouter.Use(middleware.CheckClientTokenValidation())
		authRouter.GET(urlconst.GetKeycloakUsers, controller.Controller)
		authRouter.GET(urlconst.GetKeycloakUserById, controller.Controller)
		authRouter.POST(urlconst.PostKeycloakUser, controller.Controller)
		authRouter.PUT(urlconst.PutKeycloakUsers, controller.Controller)
		authRouter.DELETE(urlconst.DeleteKeycloakUser, controller.Controller)
	}
}
