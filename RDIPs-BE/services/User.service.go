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
var ADMIN_KEYCLOAK_BASE_URL = os.Getenv("KEYCLOAK_BASE_URL")
var ADMIN_KEYCLOAK_REALM_NAME = os.Getenv("KEYCLOAK_REALM_NAME")

var PostKeycloakUser = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostKeycloakUser Start")
	defer utils.Log(LogConstant.Info, "PostKeycloakUser End")

	client := gocloak.NewClient(ADMIN_KEYCLOAK_BASE_URL)
	ctx := c.Ctx.Request.Context()

	userBody := gocloak.User{}
	err := json.Unmarshal(c.Body, &userBody)

	if err == nil {
		_, err := client.CreateUser(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			userBody,
		)
		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
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

var PostKeycloakGroup = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostKeycloakGroup Start")
	defer utils.Log(LogConstant.Info, "PostKeycloakGroup End")

	client := gocloak.NewClient(ADMIN_KEYCLOAK_BASE_URL)
	ctx := c.Ctx.Request.Context()

	groupRequest := model.GroupRequest{}
	err := json.Unmarshal(c.Body, &groupRequest)

	groupBody := gocloak.Group{
		Name:       groupRequest.Name,
		Attributes: groupRequest.Attributes,
		RealmRoles: groupRequest.RealmRoles,
	}

	if err == nil {
		roles, err := getRealRoles(c, client, groupBody)

		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}

		groupID, err := client.CreateGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			groupBody,
		)

		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}

		err = client.AddRealmRoleToGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
			roles,
		)

		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}

		err = addUsersToGroup(c, client, groupRequest.UserIds, groupID)

		if err != nil {
			client.DeleteGroup(
				ctx,
				c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
				ADMIN_KEYCLOAK_REALM_NAME,
				groupID,
			)

			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
	}

	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, err
}

var getRealRoles = func(c *commonModel.ServiceContext,
	client *gocloak.GoCloak, groupBody gocloak.Group) ([]gocloak.Role, error) {
	ctx := c.Ctx.Request.Context()
	roles := []gocloak.Role{}
	for _, role := range *groupBody.RealmRoles {
		role, err := client.GetRealmRole(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			role,
		)

		if err != nil {
			return nil, err
		}

		roles = append(roles, *role)
	}

	return roles, nil
}

var addUsersToGroup = func(c *commonModel.ServiceContext,
	client *gocloak.GoCloak, userIds *[]string, groupID string) error {
	ctx := c.Ctx.Request.Context()

	for _, userId := range *userIds {
		err := client.AddUserToGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			userId,
			groupID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
