package keycloak

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
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

type GoCloakClientStruct struct {
	GoCloakClient *gocloak.GoCloak
	client_id     string
	client_secret string
}

var keycloakPool handler.Pool

func InitKeycloakClient() error {
	factoryFn := func() (interface{}, error) {
		gocloakClient := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
		// restyClient := gocloakClient.RestyClient()
		// restyClient.SetDebug(true)
		// restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		ctx := context.Background()
		jwt, err := gocloakClient.LoginAdmin(ctx, os.Getenv("KEYCLOAK_USER"),
			os.Getenv("KEYCLOAK_PASSWORD"),
			ADMIN_KEYCLOAK_REALM_NAME)

		if err != nil {
			return nil, err
		}

		client, err := gocloakClient.GetClientRepresentation(ctx, jwt.AccessToken, ADMIN_KEYCLOAK_REALM_NAME, CLIENT_ID)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return nil, err
		}
		return &GoCloakClientStruct{GoCloakClient: gocloakClient, client_id: *client.ID, client_secret: *client.Secret}, err
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

	err := keycloakPool.FillPool(poolData)

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
	authURL := ADMIN_KEYCLOAK_BASE_URL + "/realms/" + ADMIN_KEYCLOAK_REALM_NAME +
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
		"client_id":     CLIENT_ID,
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

func RefreshAccessTokem(ctx context.Context, refreshToken string) (jwt *gocloak.JWT, err error) {
	goCloakObj, err := GetGocloakObj()
	gkClient := goCloakObj.GoCloakClient

	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	if refreshToken != "" {
		jwt, err := gkClient.RefreshToken(ctx, refreshToken, CLIENT_ID, goCloakObj.client_secret, ADMIN_KEYCLOAK_REALM_NAME)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return nil, err
		}
		return jwt, nil
	}
	return nil, fmt.Errorf("can't have access_token")
}
