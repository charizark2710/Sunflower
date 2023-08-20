package validators

import (
	"RDIPs-BE/model"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type HistoryValidator struct {
}

func (HistoryValidator) Post(body []byte) error {
	var validate = validator.New()
	var history model.History
	err := json.Unmarshal(body, &history)
	if err != nil {
		return err
	}
	if err := validate.Struct(&history); err != nil {
		return err
	}
	return nil
}

func (HistoryValidator) Put(body []byte) error {
	var validate = validator.New()
	var history model.History
	err := json.Unmarshal(body, &history)
	if err != nil {
		return err
	}
	if err := validate.Struct(&history); err != nil {
		return err
	}
	return nil
}
