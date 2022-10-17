package consts

// 这里定义的常量，一般是具有错误代码+错误说明组成，一般用于接口返回
var (
	// UserNotFound 无法查获用户
	UserNotFound                      = "user_not_found"
	CreateContactSuccess              = "create_contact_success"
	CreateContactFailed               = "create_contact_failed"
	CreateContactCannotAddOwnAcc      = "create_contact_cannot_add_own_acc"
	CreateContactSearchContactCrashed = "create_contact_search_contact_crashed"
	CreateContactRequestDuplicated    = "create_contact_request_duplicated"
	CreateGroupSuccess                = "create_group_success"
	CreateGroupFailed                 = "create_group_failed"
	AddGroupMemberSuccess             = "add_group_member_success"
	AddGroupMemberFailed              = "add_group_member_failed"
	SetGroupAdminSuccess              = "set_group_admin_success"
	SetGroupAdminFailed               = "set_group_admin_failed"
	SetGroupOwnerSuccess              = "set_group_owner_success"
	SetGroupOwnerFailed               = "set_group_owner_failed"
	GetGroupInfoSuccess               = "get_group_info_success"
	GetGroupInfoFailed                = "get_group_info_failed"
	NoGroupAuthority                  = "no_group_authority"
	MemberExistInGroup                = "member_already_in_group"
	MemberNotInGroup                  = "member_not_in_group"
	RemoveGroupMemberSuccess          = "remove_group_member_success"
	RemoveGroupMemberFailed           = "remove_group_member_failed"
	GroupDoesNotExist                 = "group_does_not_exist"
	DisbandGroupSuccess               = "disband_group_success"
	DisbandGroupFailed                = "disband_group_failed"
)
