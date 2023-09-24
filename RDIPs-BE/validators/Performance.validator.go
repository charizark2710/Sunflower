package validators

import (
	"RDIPs-BE/model"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type PerformanceValidator struct {
}

func (PerformanceValidator) Post(body []byte) error {
	var validate = validator.New()
	var performance model.Performance
	err := json.Unmarshal(body, &performance)
	if err != nil {
		return err
	}
	if err := validate.Struct(&performance); err != nil {
		return err
	}
	return nil
}

func (PerformanceValidator) Put(body []byte) error {
	var validate = validator.New()
	var performance model.Performance
	err := json.Unmarshal(body, &performance)
	if err != nil {
		return err
	}
	if err := validate.Struct(&performance); err != nil {
		return err
	}
	return nil
}
