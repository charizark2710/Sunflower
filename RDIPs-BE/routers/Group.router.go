package routers

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/controller"

	"github.com/gin-gonic/gin"
)

func GroupRouter(router *gin.Engine) {
	router.GET(urlconst.GetKeycloakGroups, controller.Controller)
	router.GET(urlconst.GetKeycloakGroupById, controller.Controller)
	router.DELETE(urlconst.DeleteKeycloakGroup, controller.Controller)
	router.POST(urlconst.PostKeycloakGroup, controller.Controller)
	router.PUT(urlconst.PutKeycloakGroup, controller.Controller)
}
