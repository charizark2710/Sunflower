package validators

import (
	"RDIPs-BE/model"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type DeviceValidator struct {
}

func (DeviceValidator) Post(body []byte) error {
	var validate = validator.New()
	var device model.Devices
	err := json.Unmarshal(body, &device)
	if err != nil {
		return err
	}
	if err := validate.Struct(&device); err != nil {
		return err
	}
	return nil
}

func (DeviceValidator) Put(body []byte) error {
	var validate = validator.New()
	var device model.Devices
	err := json.Unmarshal(body, &device)
	if err != nil {
		return err
	}
	if err := validate.Struct(&device); err != nil {
		return err
	}
	return nil
}
