package model

type GroupRequest struct {
	Name          *string              `json:"name,omitempty"`
	Attributes    *map[string][]string `json:"attributes,omitempty"`
	NewRealmRoles *[]string            `json:"newRealmRoles,omitempty"`
	OldRealmRoles *[]string            `json:"oldRealmRoles,omitempty"`
	NewUserIds    *[]string            `json:"newUserIds,omitempty"`
	OldUserIds    *[]string            `json:"oldUserIds,omitempty"`
}
