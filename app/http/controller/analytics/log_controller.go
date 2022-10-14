package analytics

import (
	"github.com/gin-gonic/gin"
)

func WriteLog(c *gin.Context) {
	c.String(200,"SUCCESS")
	return
}
