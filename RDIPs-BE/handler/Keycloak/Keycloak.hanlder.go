package keycloak

import (
	"RDIPs-BE/handler"
	"context"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

var keycloakPool handler.Pool

func InitKeycloakClient() error {
	factoryFn := func() (interface{}, error) {
		client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
		// restyClient := client.RestyClient()
		// restyClient.SetDebug(true)
		// restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		return client, nil
	}

	pingFn := func(conn interface{}) error {
		client := conn.(*gocloak.GoCloak)
		_, err := client.LoginAdmin(context.Background(), os.Getenv("KEYCLOAK_USER"),
			os.Getenv("KEYCLOAK_PASSWORD"),
			os.Getenv("KEYCLOAK_REALM_NAME"))
		return err
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
