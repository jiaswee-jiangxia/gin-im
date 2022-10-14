package model

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/helpers"
)

type LoginStruct struct {
	Id       int64  `gorm:"column:id" json:"-"`
	Username string `gorm:"column:username" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	Contact  string `gorm:"column:contact" json:"contact"`
}

func UserLogin(mobileNo string, pass string) (*LoginStruct, error) {
	var Member *LoginStruct
	//idx := helpers.ShardHash(mobileNo)
	hash := helpers.GetMD5Hash(pass)
	err := db.Table("users").
		Where("username = ?", mobileNo).
		Where("password = ?", hash).
		First(&Member).Error
	if err != nil {
		return nil, err
	}
	return Member, nil
}

type MainRegisterStruct struct {
	Id       int    `gorm:"primary_key" json:"id"`
	ShardKey int    `gorm:"column:shard_key" json:"shard_key"`
	MobileNo string `gorm:"column:mobile_no" json:"mobile_no"`
}

type RegisterStruct struct {
	BaseModel
	Username string `gorm:"column:username" json:"username"`
	Contact  string `gorm:"column:contact" json:"contact"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
}

func UserRegister(tx *gorm.DB, username string, email string, pass string, mobileNo string) (*RegisterStruct, error) {
	//idx := helpers.ShardHash(mobileNo)
	hash := helpers.GetMD5Hash(pass)
	//key, _ := strconv.Atoi(idx)
	var checkUser LoginStruct
	err := tx.Table("users").
		Where("username", username).First(&checkUser).Error
	if err != nil || checkUser.Id > 0 {
		return nil, errors.New("username_is_used")
	}

	//user := MainRegisterStruct{
	//	ShardKey: key,
	//	MobileNo: mobileNo,
	//}
	//var data LoginStruct
	//tx.Table("users").First(&data)
	//err = tx.Table("users").
	//	Create(&user).Error
	//if err != nil {
	//	return nil, err
	//}

	registrationClone := RegisterStruct{
		Username: username,
		Email:    email,
		Contact:  mobileNo,
		Password: hash,
	}
	err = tx.Table("users").Create(&registrationClone).Error
	if err != nil {
		return nil, err
	}
	return &registrationClone, nil
}
