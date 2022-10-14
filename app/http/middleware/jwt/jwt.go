package jwt

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/helpers"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var err string
		var claim *helpers.Claims

		code = 200

		// get header token
		token := c.GetHeader("Authorization")

		if token == "" { // empty token
			code = 401
			err = "invalid_token"
		} else {
			// parse token
			claim, err = helpers.ParseToken(token)
			if err != "" {
				code = 401
			}

			if claim != nil {
				c.Set("username", claim.Username)
				c.Set("mobile_no", claim.MobileNo)
				c.Set("user_id", claim.Id)
			} else {
				code = 401
			}
		}

		if code != 200 { // error return
			c.JSON(http.StatusUnauthorized, gin.H{
				"rst":  0,
				"msg":  err,
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}