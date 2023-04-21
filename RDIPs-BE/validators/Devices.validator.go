package validators

import (
	"RDIPs-BE/model"
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var PostDeviceValidator = func(c *gin.Context) error {
	var device model.Devices
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, &device); err != nil {
		return err
	}
	if err := validate.Struct(&device); err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
	return nil
}

var UpdateDeviceValidator = func(c *gin.Context) error {
	id := c.Param("id")
	var device model.Devices
	if err := validate.Var(id, "required,uuid"); err != nil {
		return err
	}
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, &device); err != nil {
		return err
	}
	if err := validate.Struct(&device); err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
	return nil
}
