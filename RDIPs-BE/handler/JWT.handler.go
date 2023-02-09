package handler

import (
	"fmt"

	commonModel "github.com/charizark2710/Sunflower/RDIPs-BE/model/common"
	"github.com/golang-jwt/jwt"
)

func SignToken(payload *commonModel.Credential, secrect string) (string, error) {
	// Test
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secrect))
}

func VerifyToken(token string, secrect string) (*commonModel.Credential, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &commonModel.Credential{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secrect), nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken.Claims.(*commonModel.Credential), nil

}
