package keycloak

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	"RDIPs-BE/utils"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var ADMIN_KEYCLOAK_BASE_URL = os.Getenv("KEYCLOAK_BASE_URL")
var ADMIN_KEYCLOAK_MASTER_HOST = ADMIN_KEYCLOAK_BASE_URL + "/admin/realms/master/"
var ADMIN_KEYCLOAK_REALM_NAME = os.Getenv("KEYCLOAK_REALM_NAME")
var CLIENT_ID = os.Getenv("KEYCLOAK_CLIENT_ID")
var REDIRECT_URI = os.Getenv("REDIRECT_URI")

// For debug local docker deploy only
var KEYCLOAK_AUTHEN_URL = os.Getenv("KEYCLOAK_AUTHEN_URL")

type GoCloakClientStruct struct {
	GoCloakClient *gocloak.GoCloak
	client_id     string
	client_secret string
	client_name   string
}

var keycloakPool handler.Pool

/*
Get client_id, client_secret from client_name
*/
func getClientData(client_name string) (*string, *string, error) {
	gocloakClient := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	ctx := context.Background()
	jwt, err := gocloakClient.LoginAdmin(ctx, os.Getenv("KEYCLOAK_ADMIN"),
		os.Getenv("KEYCLOAK_ADMIN_PASSWORD"),
		ADMIN_KEYCLOAK_REALM_NAME)

	if err != nil {
		return nil, nil, err
	}

	client, err := gocloakClient.GetClientRepresentation(ctx, jwt.AccessToken, ADMIN_KEYCLOAK_REALM_NAME, client_name)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, nil, err
	}
	var secretValue *string
	if client.Secret != nil && *client.Secret != "" {
		secretValue = client.Secret
	} else {
		secret, err := gocloakClient.RegenerateClientSecret(ctx, jwt.AccessToken, ADMIN_KEYCLOAK_REALM_NAME, *client.ID)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return nil, nil, err
		}

		secretValue = secret.Value
	}
	return client.ClientID, secretValue, nil
}

func InitKeycloakClient(client_name string) error {
	if client_name == "" {
		client_name = CLIENT_ID
	}
	client_id, client_secret, err := getClientData(client_name)
	if err != nil {
		return err
	}
	factoryFn := func() (interface{}, error) {

		gocloakClient := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
		// restyClient := gocloakClient.RestyClient()
		// restyClient.SetDebug(true)
		// restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		return &GoCloakClientStruct{GoCloakClient: gocloakClient, client_id: *client_id, client_secret: *client_secret, client_name: client_name}, err
	}

	pingFn := func(conn interface{}) error {
		return nil
	}

	closeFn := func(conn interface{}) error {
		return nil
	}

	poolData := handler.PoolData{
		FactoryFn: factoryFn,
		CloseFn:   closeFn,
		PingFn:    pingFn,
	}

	err = keycloakPool.FillPool(poolData)

	return err
}

func GetGocloakObj() (*GoCloakClientStruct, error) {
	keycloakObj, err := keycloakPool.Get()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}

	return keycloakObj.(*GoCloakClientStruct), nil
}

func CreateUser(ctx context.Context, clientKey string, userBody gocloak.User) (string, error) {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return "", err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient
	userId, userRrr := client.CreateUser(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		userBody,
	)
	return userId, userRrr
}

func GetUsers(ctx context.Context, clientKey string, params gocloak.GetUsersParams) ([]*gocloak.User, error) {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient
	users, userRrr := client.GetUsers(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		params,
	)
	return users, userRrr
}

func GetUserByID(ctx context.Context, clientKey string, id string) (*gocloak.User, error) {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient
	user, userRrr := client.GetUserByID(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		id,
	)
	return user, userRrr
}

func GetLoginScreen() (string, string, error) {
	codeVerifier := oauth2.GenerateVerifier()
	codeChallenge := oauth2.S256ChallengeFromVerifier(codeVerifier)
	codeVerifyMethod := "S256"
	kcEndpoint := ADMIN_KEYCLOAK_BASE_URL
	if KEYCLOAK_AUTHEN_URL != "" {
		kcEndpoint = KEYCLOAK_AUTHEN_URL
	}
	authURL := kcEndpoint + "/realms/" + ADMIN_KEYCLOAK_REALM_NAME +
		"/protocol/openid-connect/auth?response_type=code" +
		"&client_id=" + CLIENT_ID +
		"&redirect_uri=" + url.PathEscape(REDIRECT_URI) +
		"&scope=" + "openid" +
		"&state=" + uuid.NewString() +
		"&code_challenge=" + codeChallenge +
		"&code_challenge_method=" + codeVerifyMethod

	return authURL, codeVerifier, nil
}

func GetTokenObject(ctx context.Context, code, codeVerify string) (map[string]interface{}, error) {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	defer keycloakPool.Release(goCloakObj)
	authorizationToken := base64.RawURLEncoding.EncodeToString([]byte(goCloakObj.client_id + ":" + goCloakObj.client_secret))

	// TODO: http client should be in another handler
	httpClient := goCloakObj.GoCloakClient.RestyClient().R()

	endpoint := ADMIN_KEYCLOAK_BASE_URL + "/realms/" + ADMIN_KEYCLOAK_REALM_NAME + "/protocol/openid-connect/token"
	var result map[string]interface{}
	resp, err := httpClient.SetFormData(map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  REDIRECT_URI,
		"code_verifier": codeVerify,
		"client_id":     goCloakObj.client_name,
	}).SetHeader("Authorization", "Basic "+authorizationToken).SetResult(&result).Post(endpoint)

	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}

	if resp.IsError() {
		var msg string
		if e := resp.Error(); e != nil {
			msg = fmt.Sprintf("%s: %s: %s", resp.Status(), e, string(resp.Body()))
		} else {
			msg = resp.Status()
		}
		utils.Log(LogConstant.Error, msg)
		return nil, fmt.Errorf(msg)
	}

	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	return result, nil
}

func RefreshAccessToken(ctx context.Context, refreshToken string) (jwt *gocloak.JWT, err error) {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	gkClient := goCloakObj.GoCloakClient
	defer keycloakPool.Release(goCloakObj)
	if refreshToken != "" {
		jwt, err := gkClient.RefreshToken(ctx, refreshToken, goCloakObj.client_name, goCloakObj.client_secret, ADMIN_KEYCLOAK_REALM_NAME)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return nil, err
		}
		return jwt, nil
	}
	return nil, fmt.Errorf("can't have access_token")
}

func PutKeycloakUser(ctx context.Context, clientKey string, id string, userBody gocloak.User) error {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	defer keycloakPool.Release(goCloakObj)
	gkClient := goCloakObj.GoCloakClient
	user, userRrr := gkClient.GetUserByID(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		id,
	)
	if userRrr == nil && user != nil {
		user.Username = userBody.Username
		user.Enabled = userBody.Enabled
		user.FirstName = userBody.FirstName
		user.LastName = userBody.LastName
		user.Email = userBody.Email
		user.EmailVerified = userBody.EmailVerified
		updatedUserRrr := gkClient.UpdateUser(
			ctx,
			clientKey,
			ADMIN_KEYCLOAK_REALM_NAME,
			*user,
		)
		return updatedUserRrr
	}
	return userRrr
}

func AddRealmRoleToUser(ctx context.Context, clientKey string, id string, roles []gocloak.Role) error {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	gkClient := goCloakObj.GoCloakClient
	defer keycloakPool.Release(goCloakObj)

	addRoleErr := gkClient.AddRealmRoleToUser(ctx, clientKey, ADMIN_KEYCLOAK_REALM_NAME, id, roles)

	return addRoleErr
}

func DeleteKeycloakUser(ctx context.Context, clientKey string, id string) error {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	defer keycloakPool.Release(goCloakObj)
	gkClient := goCloakObj.GoCloakClient
	user, userRrr := gkClient.GetUserByID(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		id,
	)
	if userRrr == nil && user != nil {
		isEnabled := false
		user.Enabled = &isEnabled

		updatedUserRrr := gkClient.UpdateUser(
			ctx,
			clientKey,
			ADMIN_KEYCLOAK_REALM_NAME,
			*user,
		)
		return updatedUserRrr
	}
	return userRrr
}

func GetGroups(ctx context.Context, clientKey string, params gocloak.GetGroupsParams) ([]*gocloak.Group, error) {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient
	groups, err := client.GetGroups(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		params,
	)
	return groups, err
}

func GetGroupById(ctx context.Context, clientKey string, id string) (*gocloak.Group, error) {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient
	group, err := client.GetGroup(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		id,
	)
	return group, err
}

func DeleteGroup(ctx context.Context, clientKey string, id string) error {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient
	err = client.DeleteGroup(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		id,
	)
	return err
}

func CreateGroup(ctx context.Context, clientKey string,
	groupRequest model.GroupRequest, groupBody gocloak.Group) error {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient
	roles, err := getRealmRoles(ctx, clientKey, client, groupBody)
	if err != nil {
		return err
	}

	groupID, err := client.CreateGroup(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		groupBody,
	)
	if err != nil {
		return err
	}

	if len(roles) > 0 {
		err = client.AddRealmRoleToGroup(
			ctx,
			clientKey,
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
			roles,
		)
		if err != nil {
			return err
		}
	}

	var newUserIds = groupRequest.NewUserIds
	if newUserIds != nil && len(*newUserIds) > 0 {
		err = addUsersToGroup(ctx, clientKey, client, newUserIds, groupID)
	}
	return err
}

func EditGroup(ctx context.Context, clientKey string, groupID string,
	groupRequest model.GroupRequest) error {
	goCloakObj, err := GetGocloakObj()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	defer keycloakPool.Release(goCloakObj)
	client := goCloakObj.GoCloakClient

	group, err := client.GetGroup(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		groupID,
	)
	if err != nil {
		return err
	}

	group.Name = groupRequest.Name
	group.Attributes = groupRequest.Attributes
	err = client.UpdateGroup(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		*group,
	)
	if err != nil {
		return err
	}

	group.RealmRoles = groupRequest.OldRealmRoles
	oldRoles, err := getRealmRoles(ctx, clientKey, client, *group)
	if err != nil {
		return err
	}
	if len(oldRoles) > 0 {
		err = client.DeleteRealmRoleFromGroup(
			ctx,
			clientKey,
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
			oldRoles,
		)
		if err != nil {
			return err
		}
	}

	group.RealmRoles = groupRequest.NewRealmRoles
	newRoles, err := getRealmRoles(ctx, clientKey, client, *group)
	if err != nil {
		return err
	}
	if len(newRoles) > 0 {
		err = client.AddRealmRoleToGroup(
			ctx,
			clientKey,
			ADMIN_KEYCLOAK_REALM_NAME,
			groupID,
			newRoles,
		)
		if err != nil {
			return err
		}
	}

	var oldUserIds = groupRequest.OldUserIds
	if oldUserIds != nil && len(*oldUserIds) > 0 {
		err = deleteUsersFromGroup(ctx, clientKey, client, oldUserIds, groupID)
		if err != nil {
			return err
		}
	}

	var newUserIds = groupRequest.NewUserIds
	if newUserIds != nil && len(*newUserIds) > 0 {
		err = addUsersToGroup(ctx, clientKey, client, newUserIds, groupID)
	}

	return err
}

var getRealmRoles = func(ctx context.Context, clientKey string,
	client *gocloak.GoCloak, groupBody gocloak.Group) ([]gocloak.Role, error) {
	roles := []gocloak.Role{}
	var realmRoles = groupBody.RealmRoles
	if realmRoles != nil && len(*realmRoles) > 0 {
		for _, role := range *realmRoles {
			role, err := client.GetRealmRole(
				ctx,
				clientKey,
				ADMIN_KEYCLOAK_REALM_NAME,
				role,
			)

			if err != nil {
				return nil, err
			}

			roles = append(roles, *role)
		}
	}

	return roles, nil
}

var addUsersToGroup = func(ctx context.Context, clientKey string,
	client *gocloak.GoCloak, userIds *[]string, groupID string) error {

	for _, userId := range *userIds {
		err := client.AddUserToGroup(
			ctx,
			clientKey,
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

var deleteUsersFromGroup = func(ctx context.Context, clientKey string,
	client *gocloak.GoCloak, userIds *[]string, groupID string) error {

	for _, userId := range *userIds {
		err := client.DeleteUserFromGroup(
			ctx,
			clientKey,
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
