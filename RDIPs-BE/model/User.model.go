package model

type SysUser interface {
}

type TokenConnectResponse struct {
	Token     string `json:"access_token"`
	ExpiredIn int    `json:"expires_in"`
}

type KeycloakUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	Enabled       bool   `json:"enabled"`
}

type KeycloakUserRequest struct {
	Id            string                    `json:"id"`
	Username      string                    `json:"username"`
	FirstName     string                    `json:"firstName"`
	LastName      string                    `json:"lastName"`
	Email         string                    `json:"email"`
	EmailVerified bool                      `json:"emailVerified"`
	Enabled       bool                      `json:"enabled"`
	Credentials   []KeycloakCredentialModel `json:"credentials"`
}

type KeycloakCredentialModel struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type KeycloakRole struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
