package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	keycloak "RDIPs-BE/handler/Keycloak"
	"RDIPs-BE/middleware"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
)

var ADMIN_KEYCLOAK_MASTER_HOST = os.Getenv("KEYCLOAK_BASE_URL") + "/admin/realms/master/"
var ADMIN_KEYCLOAK_BASE_URL = os.Getenv("KEYCLOAK_BASE_URL")
var ADMIN_KEYCLOAK_REALM_NAME = os.Getenv("KEYCLOAK_REALM_NAME")
var APP_HOST = os.Getenv("APP_HOST")
var REACT_APP_API_URL = os.Getenv("REACT_APP_API_URL")

var PostKeycloakUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostKeycloakUser Start")
	defer utils.Log(LogConstant.Info, "PostKeycloakUser End")

	userBody := gocloak.User{}
	err := json.Unmarshal(c.Body, &userBody)

	if err == nil {
		_, userErr := keycloak.CreateUser(c.Ctx, c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY), userBody)
		if userErr != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, userErr
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: "Create User successfully"}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var GetKeycloakUsers = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakUsers Start")
	users, err := keycloak.GetUsers(c.Ctx, c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY), gocloak.GetUsersParams{})
	if err == nil {
		utils.Log(LogConstant.Info, "GetKeycloakUsers End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: users}, nil

	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var GetKeycloakUserById = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakUserById Start")
	userId := c.Param("id")
	user, err := keycloak.GetUserByID(c.Ctx, c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY), userId)
	if err == nil {
		utils.Log(LogConstant.Info, "GetKeycloakUserById End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: user}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

/*
* This Function will send back url of login page of keycloak
 */
var GetLoginScreen = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Debug, "GetLoginScreen Start")
	loginPage, codeVerify, err := keycloak.GetLoginScreen()

	if err == nil {
		c.Ctx.SetCookie("code", codeVerify, 5*60, "/", APP_HOST, true, true)
		c.Ctx.Header("Location", loginPage)
		return commonModel.ResponseTemplate{HttpCode: http.StatusFound, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

/*
* This Function will verify the code and return access_token
 */
var Callback = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Debug, "Callback Start")
	codeVerifier, err := c.Ctx.Cookie("code")
	c.Ctx.SetCookie("code", "", -1, "/", APP_HOST, true, true)
	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	if codeVerifier == "" {
		return commonModel.ResponseTemplate{HttpCode: 403, Data: "Unauthenticated"}, err
	}
	res, err := keycloak.GetTokenObject(c.Ctx, c.Query("code"), codeVerifier)
	if err == nil {
		unexpectedErr := fmt.Errorf("something went wrong")
		accessToken, ok := res["access_token"].(string)
		if !ok {
			utils.Log(LogConstant.Error, unexpectedErr)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, unexpectedErr
		}
		claimToken, ok := handler.ClaimsToken(accessToken)
		if !ok {
			utils.Log(LogConstant.Error, unexpectedErr)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, unexpectedErr
		}
		refreshToken, ok := res["refresh_token"].(string)
		if !ok {
			utils.Log(LogConstant.Error, unexpectedErr)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, unexpectedErr
		}
		sub, ok := claimToken["sub"].(string)
		if !ok {
			utils.Log(LogConstant.Error, unexpectedErr)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, unexpectedErr
		}
		commonModel.CacheSrv.Add(&memcache.Item{
			Key:        sub,
			Value:      []byte(refreshToken),
			Expiration: 30 * 60,
		})
		if !ok {
			utils.Log(LogConstant.Error, unexpectedErr)
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, unexpectedErr
		}
		c.Ctx.SetCookie("access_token", accessToken, 10*60, "/", APP_HOST, true, true)
		c.Ctx.SetCookie("token", uuid.NewString(), 30*60, "/", APP_HOST, true, false)

		c.Ctx.Header("Location", REACT_APP_API_URL)

		c.Ctx.Request.URL.RawQuery = ""
		return commonModel.ResponseTemplate{HttpCode: http.StatusPermanentRedirect, Data: nil}, err
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
			url := ADMIN_KEYCLOAK_MASTER_HOST + "users/" + userId + "/role-mappings/realm"
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

var PutKeycloakUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PutKeycloakUser Start")
	defer utils.Log(LogConstant.Info, "PutKeycloakUser End")

	userId := c.Param("id")
	client := gocloak.NewClient(ADMIN_KEYCLOAK_BASE_URL)
	ctx := c.Ctx.Request.Context()

	user, err := client.GetUserByID(
		ctx,
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		ADMIN_KEYCLOAK_REALM_NAME,
		userId,
	)

	if err == nil && user != nil {
		userRequest := gocloak.User{}
		if err := json.Unmarshal(c.Body, &userRequest); err == nil {
			user.FirstName = userRequest.FirstName
			user.LastName = userRequest.LastName
			user.Email = userRequest.Email
			user.Enabled = userRequest.Enabled

			err := client.UpdateUser(
				ctx,
				c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
				ADMIN_KEYCLOAK_REALM_NAME,
				*user,
			)

			if err != nil {
				return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, nil
			}
			return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
		}
	}

	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var DeleteKeycloakUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "DeleteKeycloakUser Start")
	defer utils.Log(LogConstant.Info, "DeleteKeycloakUser End")

	userId := c.Param("id")
	client := gocloak.NewClient(ADMIN_KEYCLOAK_BASE_URL)
	ctx := c.Ctx.Request.Context()

	user, err := client.GetUserByID(
		ctx,
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		ADMIN_KEYCLOAK_REALM_NAME,
		userId,
	)

	if err == nil && user != nil {
		isEnabled := false
		user.Enabled = &isEnabled

		err := client.UpdateUser(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			*user,
		)

		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, nil
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, nil
}
