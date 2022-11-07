package api

import (
	"github.com/gin-gonic/gin"
	consts "goskeleton/app/global/response"
	"goskeleton/app/helpers"
	"goskeleton/app/utils/response"
)

type ImCredentials struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

func ImUpdateRegister(context *gin.Context) {
	var creds ImCredentials
	if err := context.ShouldBind(&creds); err != nil {
		response.ErrorParam(context, creds)
		return
	}

	query, _ := helpers.ImSignEncryption(creds)
	response.Success(context, consts.Success, query)
	return
}
