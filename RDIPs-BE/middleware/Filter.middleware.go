package middleware

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	model "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AndOp = "AND"
	OrOp  = "OR"
	Eq    = "=="
	Ne    = "!="
	Like  = "contain"
)

var logicOperator = map[string]string{AndOp: "&&", OrOp: "||"}

// Middleware for filter
// format: filterBy=encodeURL(name==abc&&type==common||app_ver!=2&&firware_ver!=1)
// Filter Logic is (name == abc && type == common) || (app_ver != 2 && firmware_ver != 1)
// Ex: filterBy=name%3D%3Dabc%26type%3D%3Dcommon%7C%7Capp_ver%21%3D2%26firware_ver%21%3D1
func SetFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			db := model.Helper.GetDb()
			filterValue := c.Query("filterBy")
			if filterValue != "" {
				if strings.Contains(filterValue, logicOperator[OrOp]) {
					db = db.Where(handleOrLogic(filterValue))
				} else {
					db = db.Where(handleAndLogic(filterValue))
				}
			}
			c.Set("DB", db)
		}
		c.Next()
	}
}

func handleOrLogic(filter string) string {
	var result string
	orStatements := strings.Split(filter, logicOperator[OrOp])
	for _, orStatement := range orStatements {
		result += handleAndLogic(orStatement)
		result += " " + OrOp + " "
	}
	return strings.TrimSuffix(result, " "+OrOp+" ")
}

func handleAndLogic(filter string) string {
	var result string
	andStatements := strings.Split(filter, logicOperator[AndOp])
	for _, andStatement := range andStatements {
		switch {
		case strings.Contains(andStatement, Eq):
			{
				// handle equal case
				result += handleEqual(andStatement)
				break
			}
		case strings.Contains(andStatement, Ne):
			{
				// handle not equal case
				result += handleNotEqual(andStatement)
				break
			}
		case strings.Contains(andStatement, Like):
			{
				// handle contain case
				result += handleContain(andStatement)
				break
			}
		default:
			{
				utils.Log(LogConstant.Warning, "Format "+andStatement+" not exist.")
				return ""
			}
		}
		result += " " + AndOp + " "
	}

	return strings.TrimSuffix(result, " "+AndOp+" ")
}

// handle equal case
func handleEqual(value string) string {
	statements := strings.Split(value, Eq)
	return fmt.Sprintf("%s='%s'", statements[0], statements[1])
}

// handle not equal case
func handleNotEqual(value string) string {
	statements := strings.Split(value, Ne)
	return fmt.Sprintf("%s<>'%s'", statements[0], statements[1])
}

// handle contain case
func handleContain(value string) string {
	statements := strings.Split(value, Like)
	return fmt.Sprintf("%sLIKE'%s'", statements[0], "%"+statements[1]+"%")

}
