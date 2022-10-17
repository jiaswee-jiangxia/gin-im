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
			response.SuccessButFail(context, err.Error(), consts.UserNotFound, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
		}
		return
	}
	userService = user_service.TokenStruct{
		Username: createContactForm.TargetUsername,
	}
	targetUser, err := userService.FindUserByUsername()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SuccessButFail(context, err.Error(), consts.UserNotFound, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
		}
		return
	}

	if user.Id == targetUser.Id {
		response.SuccessButFail(context, err.Error(), consts.CreateContactCannotAddOwnAcc, nil)
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
				response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
			}
			response.Success(context, consts.CreateContactSuccess, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactSearchContactCrashed, nil)
		}
	} else {
		response.SuccessButFail(context, consts.CreateContactRequestDuplicated, consts.CreateContactRequestDuplicated, nil)
	}
	return
}

type AcceptContactForm struct {
	TargetUsername string `json:"target_username" binding:"required"`
}

func AcceptContact(context *gin.Context) {
	var acceptContactForm AcceptContactForm
	if err := context.ShouldBindJSON(&acceptContactForm); err != nil {
		response.ErrorParam(context, acceptContactForm)
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
			response.SuccessButFail(context, err.Error(), consts.UserNotFound, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
		}
		return
	}
	userService = user_service.TokenStruct{
		Username: acceptContactForm.TargetUsername,
	}
	targetUser, err := userService.FindUserByUsername()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SuccessButFail(context, err.Error(), consts.UserNotFound, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
		}
		return
	}

	if user.Id == targetUser.Id {
		response.SuccessButFail(context, consts.CreateContactCannotAddOwnAcc, consts.CreateContactCannotAddOwnAcc, nil)
		return
	}

	contactService := contacts_service.ContactsStruct{
		UserId:   strconv.Itoa(int(targetUser.Id)),
		FriendId: strconv.Itoa(int(user.Id)),
		Status:   1,
	}

	contact, err := contactService.GetContactsByBothId()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.ContactNotFound, nil)
	}

	if contact.Status == 0 {
		aContact, err := contactService.AcceptContact()
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.AcceptContactFailed, nil)
		} else {
			response.Success(context, consts.AcceptContactSuccess, aContact.UpdatedAt)
		}
	} else {
		response.SuccessButFail(context, consts.AcceptContactFailed, consts.AcceptContactFailed, nil)
	}
	return
}
