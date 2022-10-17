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
)
