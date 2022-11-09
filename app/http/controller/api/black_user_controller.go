package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	consts "goskeleton/app/global/response"
	"goskeleton/app/service/black_users_service"
	"goskeleton/app/utils/response"
)

type BlackUserForm struct {
	BlackUserId string `form:"black_user_id" json:"black_user_id" binding:"required"`
}

func CreateBlackUser(context *gin.Context) {
	var blackUserForm BlackUserForm
	if err := context.ShouldBind(&blackUserForm); err != nil { // Get request data
		response.ErrorParam(context, blackUserForm)
		return
	}
	userId, exist := context.Get("user_id")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	userIdText := fmt.Sprintf("%v", userId)
	blackUserService := black_users_service.BlackUsersStruct{
		UserId:      userIdText,
		BlackUserId: blackUserForm.BlackUserId,
	}
	_, err := blackUserService.CreateNewBlackUser()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, nil)
	return
}

func RemoveBlackUser(context *gin.Context) {
	var blackUserForm BlackUserForm
	if err := context.ShouldBind(&blackUserForm); err != nil { // Get request data
		response.ErrorParam(context, blackUserForm)
		return
	}
	userId, exist := context.Get("user_id")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	userIdText := fmt.Sprintf("%v", userId)
	blackUserService := black_users_service.BlackUsersStruct{
		UserId:      userIdText,
		BlackUserId: blackUserForm.BlackUserId,
	}
	_, err := blackUserService.RemoveBlackUser()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, nil)
	return
}

type QueryBlackUserForm struct {
	Page  string `form:"page" json:"page" binding:"required"`
	Limit string `form:"limit" json:"limit" binding:"required"`
}

func GetQueryBlackList(context *gin.Context) {
	var queryBlackUserForm QueryBlackUserForm
	if err := context.ShouldBind(&queryBlackUserForm); err != nil { // Get request data
		response.ErrorParam(context, queryBlackUserForm)
		return
	}
	userId, exist := context.Get("user_id")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	userIdText := fmt.Sprintf("%v", userId)
	blackUserService := black_users_service.BlackUsersStruct{
		UserId: userIdText,
		Page:   queryBlackUserForm.Page,
		Limit:  queryBlackUserForm.Limit,
	}

	ul, err := blackUserService.QueryBlackUser()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, ul)
	return
}
