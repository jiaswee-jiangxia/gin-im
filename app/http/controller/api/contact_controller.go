package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	consts "goskeleton/app/global/response"
	"goskeleton/app/service/contacts_service"
	"goskeleton/app/service/user_service"
	"goskeleton/app/utils/response"
	"strconv"
)

type CreateContactForm struct {
	TargetUsername string `json:"target_username" binding:"required"`
}

func CreateContact(context *gin.Context) {
	var createContactForm CreateContactForm
	if err := context.ShouldBindJSON(&createContactForm); err != nil {
		response.ErrorParam(context, createContactForm)
		return
	}
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, "ok", nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	userService := user_service.TokenStruct{
		Username: usernameText,
	}
	user, err := userService.FindUserByUsername()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SuccessButFail(context, err.Error(), consts.UserNotFound.Error(), nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed.Error(), nil)
		}
		return
	}
	userService = user_service.TokenStruct{
		Username: createContactForm.TargetUsername,
	}
	targetUser, err := userService.FindUserByUsername()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SuccessButFail(context, err.Error(), consts.UserNotFound.Error(), nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed.Error(), nil)
		}
		return
	}

	if user.Id == targetUser.Id {
		response.SuccessButFail(context, err.Error(), consts.CreateContactCannotAddOwnAcc.Error(), nil)
		return
	}

	contactService := contacts_service.ContactsStruct{
		UserId:   strconv.Itoa(int(user.Id)),
		FriendId: strconv.Itoa(int(targetUser.Id)),
	}

	_, err = contactService.GetContactsByBothId()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			frdStatus := 0
			if targetUser.BFVerified == 0 {
				frdStatus = 1
			}
			contactService.Status = int64(frdStatus)
			_, err = contactService.CreateNewContact()
			if err != nil {
				response.SuccessButFail(context, err.Error(), consts.CreateContactFailed.Error(), nil)
			}
			response.Success(context, consts.CreateContactSuccess.Error(), nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactSearchContactCrashed.Error(), nil)
		}
	} else {
		response.SuccessButFail(context, consts.CreateContactRequestDuplicated.Error(), consts.CreateContactRequestDuplicated.Error(), nil)
	}
	return
}
