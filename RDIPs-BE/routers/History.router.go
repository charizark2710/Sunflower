package routers

import (
	urlconst "github.com/charizark2710/Sunflower/RDIPs-BE/constant/URLConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/controller"
	"github.com/gin-gonic/gin"
)

func HistoryRouter(router *gin.Engine) {
	router.POST(urlconst.PostHistory, controller.Controller)
	router.GET(urlconst.GetDetailHistory, controller.Controller)
}
