package middleware

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/handler"
	keycloak "RDIPs-BE/handler/Keycloak"
	model "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const KEYCLOAK_ACCESS_TOKEN = "KeycloakTokenClient"

var wg sync.WaitGroup

func Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == urlconst.GetLoginScreen || c.FullPath() == urlconst.Callback {
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == os.Getenv("API_KEY") {
			utils.Log(LogConstant.Debug, "API key is valid")
			c.Next()
			return
		}

		tokenStr, err := c.Cookie("access_token")

		//Check expired time, and get new token if refresh token valid
		if err != nil {
			utils.Log(LogConstant.Error, "Token is Missing")
			c.AbortWithStatusJSON(http.StatusForbidden, "Token is Missing")
			return
		} else {
			claims, ok := handler.ClaimsToken(tokenStr)
			if !ok {
				utils.Log(LogConstant.Error, "Wrong token format")
				c.AbortWithStatusJSON(http.StatusUnauthorized, "Wrong token format")
				return
			}
			sub, ok := claims["sub"].(string)
			if !ok {
				utils.Log(LogConstant.Error, "Wrong token format")
				c.AbortWithStatusJSON(http.StatusUnauthorized, "Wrong token format")
				return
			}
			if isTokenExpired(claims) {
				ip := c.ClientIP()
				refreshToken, err := model.CacheSrv.Get(ip + "#" + sub)
				if err != nil {
					utils.Log(LogConstant.Error, err)
					c.AbortWithStatusJSON(http.StatusUnauthorized, err)
					return
				}
				jwt, err := keycloak.RefreshAccessToken(c, string(refreshToken.Value))
				if err != nil {
					utils.Log(LogConstant.Error, err)
					c.AbortWithStatusJSON(500, err)
					return
				}
				tokenStr = jwt.AccessToken
				c.SetCookie("access_token", tokenStr, 30*60, "/", os.Getenv("APP_HOST"), true, true)
			}
			claims, ok = handler.ClaimsToken(tokenStr)

			if !ok {
				utils.Log(LogConstant.Error, "Unauthorized")
				c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
				return
			}
			//Check permission, and return 403
			// if !userHasPermission(claims, urlconst.URLRoles[c.Request.Method+c.FullPath()]) {
			// 	utils.Log(LogConstant.Error, "User doesn't have permission")
			// 	c.AbortWithStatusJSON(http.StatusForbidden, "User doesn't have permission")
			// 	return
			// }
			utils.Log(LogConstant.Debug, "CheckPermission End")
			c.Next()
			return
		}
	}
}

/*
In case multiple requests access at the same time, one process will call api to get access_token for KeycloakTokenClient.
Other will wait until the access_token is retrieved and store to memcache.
After one request is done, the other processes will continue.
Other processes will check data in memcache to verify token if it is valid
*/
func CheckClientTokenValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		wg.Wait()
		if isKeyCloakTokenClientExpired(c) {
			wg.Add(1)
			err := getTokenAdmin(c.Request.Context(), c)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				wg.Done()
				return
			}
		}
		c.Next()
	}
}

func getTokenAdmin(ctx context.Context, c *gin.Context) error {
	defer wg.Done()
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	token, err := client.LoginAdmin(
		ctx,
		os.Getenv("KC_BOOTSTRAP_ADMIN_USERNAME"),
		os.Getenv("KC_BOOTSTRAP_ADMIN_PASSWORD"),
		os.Getenv("KEYCLOAK_REALM_NAME"))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return err
	}
	model.CacheSrv.Add(&memcache.Item{
		Key:        KEYCLOAK_ACCESS_TOKEN,
		Value:      []byte(token.AccessToken),
		Expiration: 5 * 60,
	})
	c.Set(KEYCLOAK_ACCESS_TOKEN, token.AccessToken)
	return nil

}

func isKeyCloakTokenClientExpired(c *gin.Context) bool {
	keycloakTokenItem, err := model.CacheSrv.Get(KEYCLOAK_ACCESS_TOKEN)
	if err != nil {
		return true
	}
	c.Set(KEYCLOAK_ACCESS_TOKEN, string(keycloakTokenItem.Value))
	return false
}

func isTokenExpired(claims jwt.MapClaims) bool {
	utils.Log(LogConstant.Debug, "check token is expired start")
	exp := claims["exp"].(float64)
	expUnix := int64(exp)
	// Convert Unix timestamp to time.Time
	expTime := time.Unix(expUnix, 0)
	// Verify the token's expiration time
	return !time.Now().Before(expTime)
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
