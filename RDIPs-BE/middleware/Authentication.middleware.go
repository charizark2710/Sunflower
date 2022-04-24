package middleware

import (
	"os"

	"github.com/charizark2710/Automate-Garden/RDIPs-BE/handler"
	"github.com/charizark2710/Automate-Garden/RDIPs-BE/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Validation() gin.HandlerFunc {

	return func(c *gin.Context) {
		secret := os.Getenv("SECRECT")

		if c.GetHeader("Authorization") != "" {

			_, verifyErr := handler.VerifyToken(c.GetHeader("Authorization"), secret)
			if verifyErr != nil {
				c.AbortWithStatus(401)
			}
			c.Next()
		} else {
			tokenID, uuidErr := uuid.NewRandom()
			if uuidErr != nil {
				c.AbortWithStatus(500)
				return
			}
			token, signError := handler.SignToken(&model.Credential{UserName: "test", OtherInfo: "", ID: tokenID}, secret)
			if signError != nil {
				c.AbortWithStatus(500)
				return
			}

			c.SetCookie("token", token, 0, "/", os.Getenv("HOST"), false, true)
			c.Next()
		}
	}
}
