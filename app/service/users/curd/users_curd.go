package curd

import (
	"goskeleton/app/model/base"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateUserCurdFactory() *UsersCurd {

	return &UsersCurd{}
}

type UsersCurd struct {
}

func (u *UsersCurd) Register(userName, pass, userIp string) bool {
	pass = md5_encrypt.Base64Md5(pass) // 预先处理密码加密，然后存储在数据库
	return base.CreateUserFactory("").Register(userName, pass, userIp)
}

func (u *UsersCurd) Store(name string, pass string, realName string, phone string, remark string) bool {

	pass = md5_encrypt.Base64Md5(pass) // 预先处理密码加密，然后存储在数据库
	return base.CreateUserFactory("").Store(name, pass, realName, phone, remark)
}

func (u *UsersCurd) Update(id int, name string, pass string, realName string, phone string, remark string, clientIp string) bool {
	//预先处理密码加密等操作，然后进行更新
	pass = md5_encrypt.Base64Md5(pass) // 预先处理密码加密，然后存储在数据库
	return base.CreateUserFactory("").Update(id, name, pass, realName, phone, remark, clientIp)
}
