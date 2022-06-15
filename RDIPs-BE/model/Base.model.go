package model

import "fmt"

var DbHelper *helper = &helper{}

func (credential *Credential) Valid() error {
	fmt.Print(credential)
	return nil
}
