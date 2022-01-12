package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Credential struct {
	UserName  string      `json:"userName"`
	OtherInfo interface{} `json:"otherInfo"`
	ID        uuid.UUID   `json:"id"`
}

func (credential *Credential) Valid() error {
	fmt.Print(credential)
	return nil
}
