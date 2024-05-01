package keycloak

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/utils"
	"context"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

var ADMIN_KEYCLOAK_HOST = os.Getenv("KEYCLOAK_BASE_URL") + "/admin/realms/master/"
var ADMIN_KEYCLOAK_BASE_URL = os.Getenv("KEYCLOAK_BASE_URL")
var ADMIN_KEYCLOAK_REALM_NAME = os.Getenv("KEYCLOAK_REALM_NAME")

type KeycloakClientStruct struct {
	Client *gocloak.GoCloak
	JWT    *gocloak.JWT
}

var keycloakPool handler.Pool

func InitKeycloakClient() error {
	factoryFn := func() (interface{}, error) {
		client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
		// restyClient := client.RestyClient()
		// restyClient.SetDebug(true)
		// restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		jwt, err := client.LoginAdmin(context.Background(), os.Getenv("KEYCLOAK_USER"),
			os.Getenv("KEYCLOAK_PASSWORD"),
			os.Getenv("KEYCLOAK_REALM_NAME"))
		return &KeycloakClientStruct{Client: client, JWT: jwt}, err
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

func GetKeycloakClient() (*KeycloakClientStruct, error) {
	keycloakObj, err := keycloakPool.Get()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}

	return keycloakObj.(*KeycloakClientStruct), nil
}

func CreateUser(ctx context.Context, clientKey string, userBody gocloak.User) (string, error) {
	clientObj, err := GetKeycloakClient()
	if err != nil {
		return "", err
	}
	defer keycloakPool.Release(nil)
	client := clientObj.Client
	userId, userRrr := client.CreateUser(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		userBody,
	)
	return userId, userRrr
}

func GetUsers(ctx context.Context, clientKey string, params gocloak.GetUsersParams) ([]*gocloak.User, error) {
	clientObj, err := GetKeycloakClient()
	if err != nil {
		return nil, err
	}
	defer keycloakPool.Release(nil)
	client := clientObj.Client
	users, userRrr := client.GetUsers(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		params,
	)
	return users, userRrr
}

func GetUserByID(ctx context.Context, clientKey string, id string) (*gocloak.User, error) {
	clientObj, err := GetKeycloakClient()
	if err != nil {
		return nil, err
	}
	defer keycloakPool.Release(nil)
	client := clientObj.Client
	user, userRrr := client.GetUserByID(
		ctx,
		clientKey,
		ADMIN_KEYCLOAK_REALM_NAME,
		id,
	)
	return user, userRrr
}
