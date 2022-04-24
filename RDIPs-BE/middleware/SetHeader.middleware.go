package middleware

import "github.com/gin-gonic/gin"

func SetHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "X-Requested-With,content-type")
		c.Next()
	}
}
