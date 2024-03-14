package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/middleware"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

var ADMIN_KEYCLOAK_HOST = os.Getenv("KEYCLOAK_BASE_URL") + "/admin/realms/master/"

var PostKeycloakUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostUser Start")
	userBody := model.KeycloakUserRequest{}
	err := json.Unmarshal(c.Body, &userBody)
	if err == nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(userBody); err == nil {
			url := ADMIN_KEYCLOAK_HOST + "users"
			req, _ := http.NewRequest(http.MethodPost, url, &buf)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusCreated {
				return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
			}
		}
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err

}

var GetKeycloakUsers = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakUsers Start")
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	users, err := client.GetUsers(
		c.Ctx.Request.Context(),
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		os.Getenv("KEYCLOAK_REALM_NAME"),
		gocloak.GetUsersParams{},
	)
	if err == nil {
		utils.Log(LogConstant.Info, "GetKeycloakUsers End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: users}, nil

	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var GetKeycloakUserById = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakUserById Start")
	userId := c.Param("id")
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	user, err := client.GetUserByID(
		c.Ctx.Request.Context(),
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		os.Getenv("KEYCLOAK_REALM_NAME"),
		userId,
	)
	if err == nil {
		utils.Log(LogConstant.Info, "GetKeycloakUserById End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: user}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var Login = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Debug, "Login start")
	loginRequest := model.LoginByKeycloakRequest{}
	err := json.Unmarshal(c.Body, &loginRequest)
	if err == nil {
		payload := strings.NewReader(
			"client_id=" + os.Getenv("KEYCLOAK_CLIENT_ID") +
				"&username=" + loginRequest.Username +
				"&password=" + loginRequest.Password +
				"&grant_type=password")
		url := os.Getenv("KEYCLOAK_BASE_URL") + "/realms/master/protocol/openid-connect/token"
		req, _ := http.NewRequest(http.MethodPost, url, payload)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var body, _ = io.ReadAll(resp.Body)
			var tokenConnectResponse model.TokenConnectResponse
			if err := json.Unmarshal(body, &tokenConnectResponse); err == nil {
				utils.Log(LogConstant.Debug, "Login end")
				return commonModel.ResponseTemplate{HttpCode: 200, Data: tokenConnectResponse}, nil
			}

			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}

		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var PostRoleToUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Debug, "PostRoleToUser start")
	userId := c.Param("userId")
	roleRequest := model.KeycloakRole{}
	err := json.Unmarshal(c.Body, &roleRequest)
	if err == nil {

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(roleRequest); err == nil {
			url := ADMIN_KEYCLOAK_HOST + "users/" + userId + "/role-mappings/realm"
			req, _ := http.NewRequest(http.MethodPost, url, &buf)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusCreated {
				return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
			}
		}
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}
