package handler

import (
	"../model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func SignToken(payload *model.Credential, secrect string) (string, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		credential := &model.Credential{UserName: "test", OtherInfo: "abcdef", ID: tokenID}

		// Test
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, credential)
		return token.SignedString([]byte(secrect))
	}

}
