package middleware

import (
	model "RDIPs-BE/model/common"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	ASC        = "ASC"
	DESC       = "DESC"
	FIELD_NAME = "fieldName"
	DIRECTION  = "direction"
)

// Middleware for sort
// Sort by ASC/DESC for input fieldName
// format: sortBy=encodeURL(fieldName=<fieldName>&&direction=ASC/DESC)
// Ex: fieldName=name&direction=ASC
// Encode: sortBy=fieldName%3Dname%26direction%3DASC
func SetSort() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			db, exist := c.Get(string(model.Postgres))
			if !exist {
				db, _ = model.Factory.GetGormDB()
			}
			sortByValue := c.Query("sortBy")
			if sortByValue != "" {
				sortValues := strings.Split(sortByValue, logicOperator[AndOp])
				fieldName := ""
				direction := ""
				for _, v := range sortValues {
					statements := strings.Split(v, "=")
					if statements[0] == FIELD_NAME {
						fieldName = statements[1]
					}
					if statements[0] == DIRECTION {
						direction = statements[1]
					}
				}
				db = db.(*gorm.DB).Order(fieldName + " " + direction)
			}
			c.Set(string(model.Postgres), db)
		}
		c.Next()
	}
}
