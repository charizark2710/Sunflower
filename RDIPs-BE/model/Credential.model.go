package model

import (
	"github.com/google/uuid"
)

type Credential struct {
	UserName  string      `json:"userName"`
	OtherInfo interface{} `json:"otherInfo"`
	ID        uuid.UUID   `json:"id"`
}
