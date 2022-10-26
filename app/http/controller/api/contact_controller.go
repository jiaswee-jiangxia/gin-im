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
	if err := context.ShouldBind(&createContactForm); err != nil {
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
		response.SuccessButFail(context, consts.CreateContactCannotAddOwnAcc, consts.CreateContactCannotAddOwnAcc, nil)
		return
	}

	contactService := contacts_service.ContactsStruct{
		UserId:   strconv.Itoa(int(user.Id)),
		FriendId: strconv.Itoa(int(targetUser.Id)),
	}
	_, err = contactService.GetContactsByBothId()

	contactService2 := contacts_service.ContactsStruct{
		UserId:   strconv.Itoa(int(targetUser.Id)),
		FriendId: strconv.Itoa(int(user.Id)),
	}
	contact2, err2 := contactService2.GetContactsByBothId()
	frdStatus := 1
	if err2 == gorm.ErrRecordNotFound {
		frdStatus = 0
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			contactService.Status = int64(frdStatus)
			_, err = contactService.CreateNewContact()
			if err != nil {
				response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
			}
			if contact2 != nil && contact2.UserId > 0 {
				contactService2.Status = int64(frdStatus)
				_, err = contactService2.UpdateContact()
				if err != nil {
					response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
					return
				}
			}
			response.Success(context, consts.CreateContactSuccess, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateContactSearchContactCrashed, nil)
		}
	}
	//else {
	//	if contact.Status > -1 {
	//		response.SuccessButFail(context, consts.CreateContactRequestDuplicated, consts.CreateContactRequestDuplicated, nil)
	//		return
	//	} else {
	//		contactService.Status = int64(frdStatus)
	//		_, err = contactService.UpdateContact()
	//		if err != nil {
	//			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
	//		}
	//		contactService.UserId = strconv.Itoa(int(targetUser.Id))
	//		contactService.FriendId = strconv.Itoa(int(user.Id))
	//		_, err = contactService.UpdateContact()
	//		if err != nil {
	//			response.SuccessButFail(context, err.Error(), consts.CreateContactFailed, nil)
	//		}
	//	}
	//	response.Success(context, consts.CreateContactSuccess, nil)
	//}
	return
}

type AcceptContactForm struct {
	TargetUsername string `json:"target_username" binding:"required"`
}

func AcceptContact(context *gin.Context) {
	var acceptContactForm AcceptContactForm
	if err := context.ShouldBind(&acceptContactForm); err != nil {
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

type RemoveContactForm struct {
	TargetUsername string `json:"target_username" binding:"required"`
}

func RemoveContact(context *gin.Context) {
	var removeContactForm RemoveContactForm
	if err := context.ShouldBind(&removeContactForm); err != nil {
		response.ErrorParam(context, removeContactForm)
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
			response.SuccessButFail(context, err.Error(), consts.RemoveGroupMemberFailed, nil)
		}
		return
	}
	userService = user_service.TokenStruct{
		Username: removeContactForm.TargetUsername,
	}
	targetUser, err := userService.FindUserByUsername()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SuccessButFail(context, err.Error(), consts.UserNotFound, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.RemoveGroupMemberFailed, nil)
		}
		return
	}

	if user.Id == targetUser.Id {
		response.SuccessButFail(context, consts.ContactCannotRemoveOwnAcc, consts.ContactCannotRemoveOwnAcc, nil)
		return
	}

	contactService := contacts_service.ContactsStruct{
		UserId:   strconv.Itoa(int(user.Id)),
		FriendId: strconv.Itoa(int(targetUser.Id)),
		Status:   -1,
	}

	contact, err := contactService.GetContactsByBothId()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.ContactNotFound, nil)
	}

	if contact.Status >= 0 {
		_, err = contactService.UpdateContact()
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.RemoveContactFailed, nil)
		}
		contactService.UserId = strconv.Itoa(int(targetUser.Id))
		contactService.FriendId = strconv.Itoa(int(user.Id))
		aFrdContact, err := contactService.UpdateContact()
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.RemoveContactFailed, nil)
		}
		response.Success(context, consts.RemoveContactSuccess, aFrdContact.UpdatedAt)
	} else {
		response.SuccessButFail(context, consts.RemoveContactFailed, consts.RemoveContactFailed, nil)
	}
	return
}

type CreateGroupingForm struct {
	TargetUsername string `json:"target_username" binding:"required"`
	GroupName      string `json:"group_name" binding:"required"`
}

func CreateGrouping(context *gin.Context) {
	var createGroupingForm CreateGroupingForm
	if err := context.ShouldBind(&createGroupingForm); err != nil {
		response.ErrorParam(context, createGroupingForm)
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
			response.SuccessButFail(context, err.Error(), consts.CreateGroupingFailed, nil)
		}
		return
	}
	userService = user_service.TokenStruct{
		Username: createGroupingForm.TargetUsername,
	}
	targetUser, err := userService.FindUserByUsername()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SuccessButFail(context, err.Error(), consts.UserNotFound, nil)
		} else {
			response.SuccessButFail(context, err.Error(), consts.CreateGroupingFailed, nil)
		}
		return
	}

	if user.Id == targetUser.Id {
		response.SuccessButFail(context, consts.CannotCreateGroupingOwnAcc, consts.CannotCreateGroupingOwnAcc, nil)
		return
	}

	contactService := contacts_service.ContactsStruct{
		UserId:   strconv.Itoa(int(user.Id)),
		FriendId: strconv.Itoa(int(targetUser.Id)),
		Grouping: createGroupingForm.GroupName,
	}

	contact, err := contactService.GetContactsByBothId()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.ContactNotFound, nil)
	}

	if contact.Status >= 0 {
		uContact, err := contactService.UpdateContactGrouping()
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.CannotCreateGroupingOwnAcc, nil)
		}
		response.Success(context, consts.CreateGroupingSuccess, uContact.UpdatedAt)
	} else {
		response.SuccessButFail(context, consts.CannotCreateGroupingOwnAcc, consts.CannotCreateGroupingOwnAcc, nil)
	}
	return
}

type SearchContactForm struct {
	TargetUsername string `form:"target_username" json:"target_username" binding:"required"`
}

func SearchContact(context *gin.Context) {
	var searchContactForm SearchContactForm
	if err := context.ShouldBind(&searchContactForm); err != nil {
		response.ErrorParam(context, searchContactForm)
		return
	}
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, "ok", nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	if usernameText == searchContactForm.TargetUsername {
		response.SuccessButFail(context, consts.UserNotFound, consts.UserNotFound, nil)
		return
	}
	userService := user_service.TokenStruct{
		Username: searchContactForm.TargetUsername,
	}
	targetUser, err := userService.FindUserByUsername()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.UserNotFound, nil)
		return
	}

	response.Success(context, consts.UserFound, map[string]string{
		"username": targetUser.Username,
		"email":    targetUser.Email,
		"contact":  targetUser.Contact,
	})
	return
}

type ContactListForm struct {
}

type ContactListItem struct {
	Username string `json:"username"`
	Status   string `json:"status"`
}

func ContactList(context *gin.Context) {
	var contactListForm ContactListForm
	var returnedContactList []ContactListItem
	if err := context.ShouldBind(&contactListForm); err != nil {
		response.ErrorParam(context, contactListForm)
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
			response.SuccessButFail(context, err.Error(), consts.CreateGroupingFailed, nil)
		}
		return
	}
	contactService := contacts_service.ContactsStruct{
		UserId: strconv.Itoa(int(user.Id)),
	}
	list, err := contactService.GetContactList()
	for _, item := range list {
		returnedContactList = append(returnedContactList, ContactListItem{
			Username: item.Username,
			Status:   strconv.Itoa(int(item.Status)),
		})
	}
	response.Success(context, consts.Success, returnedContactList)
	return
}
