package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func SignToken(payload *commonModel.Credential, secrect string) (string, error) {
	// Test
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secrect))
}

func VerifyToken(token string, secrect string) (jwt.MapClaims, error) {
	token = strings.Replace(token, "Bearer ", "", 1)

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secrect), nil
	})
	if err != nil {
		return nil, err
	}

	return jwtToken.Claims.(jwt.MapClaims), nil
}

func ClaimsToken(tokenString string) (jwt.MapClaims, bool) {
	utils.Log(LogConstant.Debug, "CheckPermission Start")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEYCLOAK_PUBLIC_KEY")), nil
	})

	utils.Log(LogConstant.Debug, "claims")
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}
