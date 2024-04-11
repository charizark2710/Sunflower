package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/middleware"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

var GetKeycloakGroups = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakGroups Start")
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	utils.Log(LogConstant.Debug, "Hello "+c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY))
	groups, err := client.GetGroups(
		c.Ctx.Request.Context(),
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		os.Getenv("KEYCLOAK_REALM_NAME"),
		gocloak.GetGroupsParams{},
	)
	if err == nil {
		utils.Log(LogConstant.Info, "GetKeycloakGroups End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: groups}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var GetKeycloakGroupById = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetKeycloakGroupById Start")
	id := c.Param("id")
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	group, err := client.GetGroup(
		c.Ctx.Request.Context(),
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		os.Getenv("KEYCLOAK_REALM_NAME"),
		id,
	)
	if err == nil {
		utils.Log(LogConstant.Info, "GetKeycloakGroupById End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: group}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var DeleteKeycloakGroup = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "DeleteKeycloakGroup Start")
	defer utils.Log(LogConstant.Info, "DeleteKeycloakGroup End")
	id := c.Param("id")
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	err := client.DeleteGroup(
		c.Ctx.Request.Context(),
		c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
		os.Getenv("KEYCLOAK_REALM_NAME"),
		id,
	)
	if err == nil {
		utils.Log(LogConstant.Info, "GetKeycloakGroupById End")
		return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}

var PostKeycloakGroup = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "PostKeycloakGroup Start")
	defer utils.Log(LogConstant.Info, "PostKeycloakGroup End")

	client := gocloak.NewClient(ADMIN_KEYCLOAK_BASE_URL)
	ctx := c.Ctx.Request.Context()

	groupRequest := model.GroupRequest{}
	if err := json.Unmarshal(c.Body, &groupRequest); err == nil {
		groupBody := gocloak.Group{
			Name:       groupRequest.Name,
			Attributes: groupRequest.Attributes,
			RealmRoles: groupRequest.NewRealmRoles,
		}

		roles, err := getRealmRoles(c, client, groupBody)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		groupID, err := client.CreateGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			groupBody,
		)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		err = client.AddRealmRoleToGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
			roles,
		)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		err = addUsersToGroup(c, client, groupRequest.NewUserIds, groupID)
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
	client := gocloak.NewClient(ADMIN_KEYCLOAK_BASE_URL)
	ctx := c.Ctx.Request.Context()
	groupRequest := model.GroupRequest{}

	if err := json.Unmarshal(c.Body, &groupRequest); err == nil {
		group, err := client.GetGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
		)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		group.Name = groupRequest.Name
		group.Attributes = groupRequest.Attributes
		err = client.UpdateGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			*group,
		)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		group.RealmRoles = groupRequest.OldRealmRoles
		oldRoles, err := getRealmRoles(c, client, *group)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		err = client.DeleteRealmRoleFromGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
			oldRoles,
		)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		group.RealmRoles = groupRequest.NewRealmRoles
		newRoles, err := getRealmRoles(c, client, *group)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		err = client.AddRealmRoleToGroup(
			ctx,
			c.Ctx.GetString(middleware.KEYCLOAK_TOKEN_CLIENT_KEY),
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
			newRoles,
		)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		err = deleteUsersFromGroup(c, client, groupRequest.OldUserIds, groupID)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}

		err = addUsersToGroup(c, client, groupRequest.NewUserIds, groupID)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return commonModel.ResponseTemplate{HttpCode: 500, Message: err.Error()}, err
		}
	}

	return commonModel.ResponseTemplate{HttpCode: 200, Data: nil}, nil
}

var getRealmRoles = func(c *commonModel.ServiceContext,
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

var deleteUsersFromGroup = func(c *commonModel.ServiceContext,
	client *gocloak.GoCloak, userIds *[]string, groupID string) error {
	ctx := c.Ctx.Request.Context()

	for _, userId := range *userIds {
		err := client.DeleteUserFromGroup(
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
