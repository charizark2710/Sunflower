package validators

import (
	"RDIPs-BE/model"
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

var PostHistoryValidator = func(c *gin.Context) error {
	var history model.History
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, &history); err != nil {
		return err
	}
	if err := validate.Struct(&history); err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
	return nil
}
