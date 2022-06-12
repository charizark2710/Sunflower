package model

import "fmt"

func (credential *Credential) Valid() error {
	fmt.Print(credential)
	return nil
}
