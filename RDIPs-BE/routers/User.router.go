package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	router.GET(urlconst.GetLoginScreen, controller.Controller)
	router.GET(urlconst.Callback, controller.Controller)
	authRouter := router.Group("")
	{
		authRouter.Use(middleware.CheckClientTokenValidation())
		authRouter.GET(urlconst.GetKeycloakUsers, controller.Controller)
		authRouter.GET(urlconst.GetKeycloakUserById, controller.Controller)
		authRouter.POST(urlconst.PostKeycloakUser, controller.Controller)
		authRouter.PUT(urlconst.PutKeycloakUsers, controller.Controller)
		authRouter.DELETE(urlconst.DeleteKeycloakUser, controller.Controller)
	}
}
