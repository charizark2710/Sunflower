package model

import (
	"fmt"
)

type Credential struct {
	UserName  string      `json:"userName"`
	OtherInfo interface{} `json:"otherInfo"`
}

func (credential *Credential) Valid() error {
	fmt.Print(credential)
	return nil
}
