package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	keycloak "RDIPs-BE/handler/Keycloak"
	"RDIPs-BE/middleware"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
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
		c.Ctx.SetCookie("access_token", accessToken, 30*60, "/", APP_HOST, true, true)
		c.Ctx.SetCookie("token", uuid.NewString(), 30*60, "/", APP_HOST, true, false)

		c.Ctx.Header("Location", REACT_APP_API_URL)

		c.Ctx.Request.URL.RawQuery = ""
		return commonModel.ResponseTemplate{HttpCode: http.StatusPermanentRedirect, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var PostRoleToUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostRoleToUser start")
	defer utils.Log(LogConstant.Info, "PostRoleToUser End")
	userId := c.Param("userId")
	rolesRequest := []gocloak.Role{}
	err := json.Unmarshal(c.Body, &rolesRequest)
	if err == nil {
		addRolesErr := keycloak.AddRealmRoleToUser(
			c.Ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			userId,
			rolesRequest)
		if addRolesErr == nil {
			return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
		}
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var PutKeycloakUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PutKeycloakUser Start")
	defer utils.Log(LogConstant.Info, "PutKeycloakUser End")
	userId := c.Param("id")
	userRequest := gocloak.User{}
	if err := json.Unmarshal(c.Body, &userRequest); err == nil {
		err := keycloak.PutKeycloakUser(
			c.Ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			userId,
			userRequest)
		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, nil
}

var DeleteKeycloakUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "DeleteKeycloakUser Start")
	defer utils.Log(LogConstant.Info, "DeleteKeycloakUser End")
	userId := c.Param("id")
	err := keycloak.DeleteKeycloakUser(
		c.Ctx,
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		userId)

	if err != nil {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
}

var GetKeycloakGroups = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakGroups Start")
	defer utils.Log(LogConstant.Info, "GetKeycloakGroups End")
	groups, err := keycloak.GetGroups(
		c.Ctx,
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		gocloak.GetGroupsParams{})
	if err == nil {
		return commonModel.ResponseTemplate{HttpCode: 200, Data: groups}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var GetKeycloakGroupById = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakGroupById Start")
	defer utils.Log(LogConstant.Info, "GetKeycloakGroupById End")
	id := c.Param("id")
	group, err := keycloak.GetGroupById(
		c.Ctx,
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		id,
	)
	if err == nil {
		return commonModel.ResponseTemplate{HttpCode: 200, Data: group}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var DeleteKeycloakGroup = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "DeleteKeycloakGroup Start")
	defer utils.Log(LogConstant.Info, "DeleteKeycloakGroup End")
	id := c.Param("id")
	err := keycloak.DeleteGroup(
		c.Ctx,
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		id,
	)
	if err == nil {
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var PostKeycloakGroup = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostKeycloakGroup Start")
	defer utils.Log(LogConstant.Info, "PostKeycloakGroup End")

	groupRequest := model.GroupRequest{}
	if err := json.Unmarshal(c.Body, &groupRequest); err == nil {
		groupBody := gocloak.Group{
			Name:       groupRequest.Name,
			Attributes: groupRequest.Attributes,
			RealmRoles: groupRequest.NewRealmRoles,
		}

		err := keycloak.CreateGroup(
			c.Ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			groupRequest,
			groupBody,
		)

		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}
	} else {
		utils.Log(LogConstant.Error, err)
		return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
	}

	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
}

var PutKeycloakGroup = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PutKeycloakGroup Start")
	defer utils.Log(LogConstant.Info, "PutKeycloakGroup End")

	groupID := c.Param("id")
	groupRequest := model.GroupRequest{}

	if err := json.Unmarshal(c.Body, &groupRequest); err == nil {
		err := keycloak.EditGroup(
			c.Ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			groupID,
			groupRequest,
		)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}
	}

	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
}
