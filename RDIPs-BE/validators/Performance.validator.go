package validators

import (
	"RDIPs-BE/model"
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

var PostPerformanceValidator = func(c *gin.Context) error {
	var performance model.Performance
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, &performance); err != nil {
		return err
	}
	if err := validate.Struct(&performance); err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
	return nil
}

var UpdatePerformanceValidator = func(c *gin.Context) error {
	id := c.Param("id")
	var performance model.Performance
	if err := validate.Var(id, "required,uuid"); err != nil {
		return err
	}
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, &performance); err != nil {
		return err
	}
	if err := validate.Struct(&performance); err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
	return nil
}
