package validators

import (
	"RDIPs-BE/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var PostDeviceValidator = func(c *gin.Context) error {
	var device model.Devices
	if err := c.BindJSON(&device); err == nil {
		if err := validate.Struct(&device); err != nil {
			return err
		}
	}
	return nil
}

var UpdateDeviceValidator = func(c *gin.Context) error {
	id := c.Param("id")
	var device model.Devices
	if err := validate.Var(id, "required,uuid"); err != nil {
		return err
	}
	if err := c.BindJSON(&device); err == nil {
		if err := validate.Struct(&device); err != nil {
			return err
		}
	}
	return nil
}
