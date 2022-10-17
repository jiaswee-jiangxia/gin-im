package api

import (
	"fmt"
	"goskeleton/app/model"
	"goskeleton/app/service/group_service"
	"goskeleton/app/utils/response"

	"github.com/gin-gonic/gin"
)

type CreateGroupStruct struct {
	GroupName string `json:"groupname"`
}

func CreateGroup(context *gin.Context) {
	userName, exist := context.Get("username")
	if !exist {
		response.ErrorParam(context, "no user")
		return
	}
	usernameText := fmt.Sprintf("%v", userName) // Get username
	var group CreateGroupStruct
	if err := context.ShouldBindJSON(&group); err != nil { // Get request data
		response.ErrorParam(context, group)
		return
	}
	groupName := group.GroupName
	groupService := group_service.GroupStruct{ // Create group service
		Name:      groupName,
		Owner:     usernameText,
		CreatedBy: usernameText,
		Disbanded: false,
	}
	newGroup, err := groupService.CreateGroup() // Create new group
	if err != nil || newGroup.Id <= 0 {
		response.SuccessButFail(context, "failed to create group", nil)
		return
	}
	response.Success(context, "ok", newGroup)
	return
}

type GroupInfoRequestStruct struct {
	GroupID int64 `json:"GroupID"`
}
type GroupInfoReplyStruct struct {
	GroupInfo  model.GroupStruct         `json:"GroupInfo"`
	MemberInfo []model.GroupMemberStruct `json:"MemberInfo"`
}

func ListGroupAdmin(context *gin.Context) {
	var group GroupInfoRequestStruct
	if err := context.ShouldBindJSON(&group); err != nil {
		response.ErrorParam(context, group)
		return
	}
	groupService := group_service.GroupStruct{ // Create group service
		BaseModel: model.BaseModel{
			Id: group.GroupID,
		},
	}
	groupInfo, err := groupService.GetGroupInfo() // Get group info
	if err != nil {
		response.SuccessButFail(context, "failed to get group info", err)
		return
	}
	adminList, err := groupService.GetGroupAdmin() // Get group admin list
	if err != nil {
		response.SuccessButFail(context, "failed to get group info", err)
		return
	}
	response.Success(context, "ok", &GroupInfoReplyStruct{
		GroupInfo:  *groupInfo,
		MemberInfo: adminList,
	})
}

func ListGroupMember(context *gin.Context) {
	var group GroupInfoRequestStruct
	if err := context.ShouldBindJSON(&group); err != nil {
		response.ErrorParam(context, group)
		return
	}
	groupService := group_service.GroupStruct{ // Create group service
		BaseModel: model.BaseModel{
			Id: group.GroupID,
		},
	}
	groupInfo, err := groupService.GetGroupInfo() // Get Group Info
	if err != nil {
		response.SuccessButFail(context, "failed to get group info", err)
		return
	}
	adminList, err := groupService.GetGroupAdmin() // Get Group Admin List
	if err != nil {
		response.SuccessButFail(context, "failed to get group info", err)
		return
	}
	response.Success(context, "ok", &GroupInfoReplyStruct{
		GroupInfo:  *groupInfo,
		MemberInfo: adminList,
	})
}

type AddGroupMemberStruct struct {
	GroupID  int64    `json:"GroupID"`
	UserList []string `json:"UserList"`
}

func AddGroupMember(context *gin.Context) {
	var req AddGroupMemberStruct
	if err := context.ShouldBindJSON(&req); err != nil {
		response.ErrorParam(context, req)
		return
	}
	groupService := group_service.GroupStruct{ // Create group service
		BaseModel: model.BaseModel{
			Id: req.GroupID,
		},
	}
	check, _ := groupService.GetGroupInfo() // Check if group exist
	if check.Id <= 0 {
		response.SuccessButFail(context, "group does not exist", "0")
		return
	}
	if !CheckGroupAuthority(context, req.GroupID) {
		response.SuccessButFail(context, "no authority", "0")
		return
	}
	memberList, err := groupService.GetGroupMember() // Check if member already in the group
	if err != nil {
		response.SuccessButFail(context, "failed to add group member", err)
		return
	}
	for _, i := range req.UserList {
		if memberExist(memberList, i) {
			response.SuccessButFail(context, "member already in group", err)
			return
		}
	}
	for _, i := range req.UserList { // Add user into the group
		err := groupService.AddGroupMember(i)
		if err != nil {
			response.SuccessButFail(context, "failed to add group member", err)
			return
		}
	}
	response.Success(context, "ok", "0")
}

type SetGroupAdminStruct struct {
	GroupID        int64    `json:"GroupID"`
	MemberUsername []string `json:"MemberUsername"`
}

func SetGroupAdmin(context *gin.Context) {
	var req SetGroupAdminStruct
	if err := context.ShouldBindJSON(&req); err != nil {
		response.ErrorParam(context, req)
		return
	}
	groupService := group_service.GroupStruct{ // Create group service
		BaseModel: model.BaseModel{
			Id: req.GroupID,
		},
	}
	check, _ := groupService.GetGroupInfo() // Check if group exist
	if check.Id <= 0 {
		response.SuccessButFail(context, "group does not exist", "0")
		return
	}
	if !CheckGroupAuthority(context, req.GroupID) { // Check if user is admin or owner
		response.SuccessButFail(context, "no authority", "0")
		return
	}
	memberList, err := groupService.GetGroupMember()
	if err != nil {
		response.SuccessButFail(context, "set admin failed", "0")
	}
	for _, v := range req.MemberUsername { // Check if target member is in group
		if !memberExist(memberList, v) {
			response.SuccessButFail(context, "member not in group", "0")
			return
		}
		err = groupService.SetGroupAdmin(v)
		if err != nil {
			response.SuccessButFail(context, "set admin failed", "0")
			return
		}
	}
	response.Success(context, "ok", "0")
	return
}

type SetGroupOwnerStruct struct {
	GroupID        int64  `json:"GroupID"`
	MemberUsername string `json:"MemberUsername"`
}

func SetGroupOwner(context *gin.Context) {
	var req SetGroupOwnerStruct
	if err := context.ShouldBindJSON(&req); err != nil {
		response.ErrorParam(context, req)
		return
	}
	groupService := group_service.GroupStruct{ // Create group service
		BaseModel: model.BaseModel{
			Id: req.GroupID,
		},
	}
	check, _ := groupService.GetGroupInfo() // Check if group exist
	if check.Id <= 0 {
		response.SuccessButFail(context, "group does not exist", "0")
		return
	}
	if !CheckGroupAuthority(context, req.GroupID) { // Check if user is admin or owner
		response.SuccessButFail(context, "no authority", "0")
		return
	}
	memberList, err := groupService.GetGroupMember()
	if err != nil {
		response.SuccessButFail(context, "set owner failed", "0")
	}
	if !memberExist(memberList, req.MemberUsername) { // Check if target member is in group
		response.SuccessButFail(context, "member not in group", "0")
		return
	}
	err = groupService.SetGroupOwner(req.MemberUsername)
	if err != nil {
		response.SuccessButFail(context, "set owner failed", "0")
		return
	}
	response.Success(context, "ok", "0")
	return
}

type RemoveGroupMemberStruct struct {
	GroupID        int64    `json:"GroupID"`
	MemberUsername []string `json:"MemberUsername"`
}

func RemoveGroupMember(context *gin.Context) {
	var req RemoveGroupMemberStruct
	if err := context.ShouldBindJSON(&req); err != nil {
		response.ErrorParam(context, req)
		return
	}
	groupService := group_service.GroupStruct{ // Create group service
		BaseModel: model.BaseModel{
			Id: req.GroupID,
		},
	}
	check, _ := groupService.GetGroupInfo() // Check if group exist
	if check.Id <= 0 {
		response.SuccessButFail(context, "group does not exist", "0")
		return
	}
	if !CheckGroupAuthority(context, req.GroupID) { // Check if user is admin or owner
		response.SuccessButFail(context, "no authority", "0")
		return
	}
	memberList, err := groupService.GetGroupMember()
	if err != nil {
		response.SuccessButFail(context, "remove member failed", "0")
	}
	for _, v := range req.MemberUsername { // Check if target member is in group
		if !memberExist(memberList, v) {
			response.SuccessButFail(context, "member not in group", "0")
			return
		}
		err = groupService.RemoveGroupMember(v)
		if err != nil {
			response.SuccessButFail(context, "remove member failed", "0")
			return
		}
	}
	response.Success(context, "ok", "0")
	return
}

type DisbandGroupStruct struct {
	GroupID int64 `json:"GroupID"`
}

func DisbandGroup(context *gin.Context) {
	var req RemoveGroupMemberStruct
	if err := context.ShouldBindJSON(&req); err != nil {
		response.ErrorParam(context, req)
		return
	}
	groupService := group_service.GroupStruct{ // Create group service
		BaseModel: model.BaseModel{
			Id: req.GroupID,
		},
	}
	check, _ := groupService.GetGroupInfo() // Check if group exist
	if check.Id <= 0 {
		response.SuccessButFail(context, "group does not exist", "0")
		return
	}
	if !CheckGroupAuthority(context, req.GroupID) { // Check if user is admin or owner
		response.SuccessButFail(context, "no authority", "0")
		return
	}
	err := groupService.DisbandGroup()
	if err != nil {
		response.SuccessButFail(context, "failed to disband group", "0")
		return
	}
	response.Success(context, "ok", "0")
	return
}

// ------------------------------------------------------------------------------------------

func memberExist(memberList []model.GroupMemberStruct, memberUsername string) bool { // Check if member exist in the list
	for _, v := range memberList {
		if v.Username == memberUsername {
			return true
		}
	}

	return false
}

func CheckGroupAuthority(context *gin.Context, groupID int64) bool {
	userName, exist := context.Get("username")
	if !exist {
		return false
	}
	usernameText := fmt.Sprintf("%v", userName) // Get username
	groupService := &group_service.GroupStruct{
		BaseModel: model.BaseModel{
			Id: groupID,
		},
	}
	info, _ := groupService.GetGroupInfo() // Check if group exist
	if info.Id <= 0 {
		return false
	}
	if info.Owner == usernameText {
		return true
	}
	adminList, err := groupService.GetGroupAdmin()
	if err != nil {
		return false
	}
	if memberExist(adminList, usernameText) {
		return true
	} else {
		return false
	}
}
