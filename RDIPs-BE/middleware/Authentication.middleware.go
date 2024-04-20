package middleware

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/handler"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const KEYCLOAK_TOKEN_CLIENT_KEY = "KeycloakTokenClient"

var (
	KeycloakTokenClient string = ""
	//Set timeout for authetication key of Weather API Login
	refreshPeriodKeycloak   = 1 * time.Minute
	lastFetchedTimeKeycloak = time.Now()
)

func Validation() gin.HandlerFunc {

	return func(c *gin.Context) {
		secret := os.Getenv("SECRECT")
		tokenStr := c.GetHeader("Authorization")
		if tokenStr != "" {
			claims, ok := handler.ClaimsToken(tokenStr)
			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			//Check expired time, and return 403
			if isTokenExpired(claims) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			//Check permission, and return 403
			if !userHasPermission(claims, urlconst.URLRoles[c.Request.Method+c.FullPath()]) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			utils.Log(LogConstant.Debug, "CheckPermission End")
			c.Next()
		} else {
			tokenID, uuidErr := uuid.NewRandom()
			if uuidErr != nil {
				c.AbortWithStatus(500)
				return
			}
			token, signError := handler.SignToken(&commonModel.Credential{UserName: "test", OtherInfo: "", ID: tokenID}, secret)
			if signError != nil {
				c.AbortWithStatus(500)
				return
			}

			c.SetCookie("token", token, 0, "/", os.Getenv("HOST"), false, true)
			c.Next()
		}
	}
}

func CheckClientTokenValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isKeyCloakTokenClientExpired() {
			ctx, cancel := context.WithTimeout(c.Request.Context(), refreshPeriodKeycloak)
			defer cancel()

			err := getTokenByClientAccount(ctx, c)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		c.Next()
	}
}

func getTokenByClientAccount(ctx context.Context, c *gin.Context) error {
	utils.Log(LogConstant.Debug, "getTokenByClientAccount is calling")

	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	token, err := client.LoginAdmin(
		context.Background(),
		os.Getenv("KEYCLOAK_USER"),
		os.Getenv("KEYCLOAK_PASSWORD"),
		os.Getenv("KEYCLOAK_REALM_NAME"))
	utils.Log(LogConstant.Info, "After login admin")

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return err
	}
	c.Set(KEYCLOAK_TOKEN_CLIENT_KEY, token.AccessToken)
	refreshPeriodKeycloak = time.Duration(time.Duration(token.ExpiresIn).Seconds())
	lastFetchedTime = time.Now()
	return nil

}

func isKeyCloakTokenClientExpired() bool {
	if KeycloakTokenClient == "" {
		return true
	}

	return time.Now().After(lastFetchedTimeKeycloak.Add(refreshPeriodKeycloak))
}

func isTokenExpired(claims jwt.MapClaims) bool {
	utils.Log(LogConstant.Debug, "check token is expired start")
	exp := claims["exp"].(float64)
	expUnix := int64(exp)
	// Convert Unix timestamp to time.Time
	expTime := time.Unix(expUnix, 0)
	// Verify the token's expiration time
	if time.Now().Before(expTime) {
		return false
	}
	return true
}

func userHasPermission(claims jwt.MapClaims, requiredPermissions []string) bool {
	utils.Log(LogConstant.Debug, "userHasPermission start")
	if permissions, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := permissions["roles"].([]interface{}); ok {
			for _, role := range roles {
				for _, permission := range requiredPermissions {
					if role.(string) == permission {
						return true
					}
				}
			}
		}
	}
	return false
}
