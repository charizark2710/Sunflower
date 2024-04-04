package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/middleware"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
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
