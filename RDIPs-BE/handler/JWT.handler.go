package handler

import (
	"fmt"

	"github.com/charizark2710/Automate-Garden/RDIPs-BE/model"
	"github.com/golang-jwt/jwt"
)

func SignToken(payload *model.Credential, secrect string) (string, error) {
	// Test
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secrect))
}

func VerifyToken(token string, secrect string) (*model.Credential, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &model.Credential{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secrect), nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken.Claims.(*model.Credential), nil

}
