package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"
	"RDIPs-BE/middleware"

	"github.com/gin-gonic/gin"
)

func KeycloakRouter(router *gin.Engine) {
	router.POST(urlconst.PostLogin, controller.Controller)

	router.Use(middleware.CheckClientTokenValidation())
	router.GET(urlconst.GetKeycloakUsers, controller.Controller)
	router.GET(urlconst.GetKeycloakUserById, controller.Controller)
	router.POST(urlconst.PostKeycloakUser, controller.Controller)
	router.PUT(urlconst.PutKeycloakUsers, controller.Controller)
	router.DELETE(urlconst.DeleteKeycloakUser, controller.Controller)
	router.POST(urlconst.PostKeycloakGroup, controller.Controller)
}
