package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

type MerchantLoginStruct struct {
	Id       int64  `gorm:"column:id" json:"-"`
	Username string `gorm:"column:username" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	MobileNo string `gorm:"column:mobile_no" json:"mobile_no"`
}

func MerchantLogin(mobileNo string, pass string) (*MerchantLoginStruct, error) {
	var Member *MerchantLoginStruct
	//idx := helpers.ShardHash(mobileNo)
	hash := sha256.Sum256([]byte(pass))
	err := db.Table("merchants").
		Where("mobile_no = ?", mobileNo).
		Where("password = ?", hex.EncodeToString(hash[:])).
		First(&Member).Error
	if err != nil {
		return nil, err
	}
	return Member, nil
}

type MerchantMainRegisterStruct struct {
	Id       int    `gorm:"primary_key" json:"id"`
	ShardKey int    `gorm:"column:shard_key" json:"shard_key"`
	MobileNo string `gorm:"column:mobile_no" json:"mobile_no"`
}

type MerchantRegisterStruct struct {
	Id                int    `gorm:"primary_key" json:"id"`
	Username          string `gorm:"column:username" json:"username"`
	MobileNo          string `gorm:"column:mobile_no" json:"mobile_no"`
	Email             string `gorm:"column:email" json:"email"`
	Password          string `gorm:"column:password" json:"password"`
	SecondaryPassword string `gorm:"column:secondary_password" json:"secondary_password"`
}

func MerchantRegister(tx *gorm.DB, username string, email string, pass string, mobileNo string) (*MerchantRegisterStruct, error) {
	//idx := helpers.ShardHash(mobileNo)
	hash := sha256.Sum256([]byte(pass))
	//key, _ := strconv.Atoi(idx)
	var checkUser MerchantLoginStruct
	err := tx.Table("merchants").
		Where("mobile_no", mobileNo).First(&checkUser).Error
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

	registrationClone := MerchantRegisterStruct{
		Username:          username,
		Email:             email,
		MobileNo:          mobileNo,
		Password:          hex.EncodeToString(hash[:]),
		SecondaryPassword: hex.EncodeToString(hash[:]),
	}
	err = tx.Table("merchants").Create(&registrationClone).Error
	if err != nil {
		return nil, err
	}
	return &registrationClone, nil
}
