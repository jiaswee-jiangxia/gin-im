package cors

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/helpers"
	"net/http"
	"os"
)

var domains = []string{
	"127.0.0.1",
}

// 允许跨域
func Next() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsOnOff := os.Getenv("CORS_ON_OFF")
		method := c.Request.Method
		origin := c.GetHeader("Origin")
		if corsOnOff == "1" {
			if helpers.Contains(domains, origin) {

				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusAccepted)
		}
		c.Next()
	}
}
