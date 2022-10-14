package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

type RiderLoginStruct struct {
	Id       int64  `gorm:"column:id" json:"-"`
	Username string `gorm:"column:username" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	MobileNo string `gorm:"column:mobile_no" json:"mobile_no"`
}

func RiderLogin(mobileNo string, pass string) (*RiderLoginStruct, error) {
	var Member *RiderLoginStruct
	//idx := helpers.ShardHash(mobileNo)
	hash := sha256.Sum256([]byte(pass))
	err := db.Table("riders").
		Where("mobile_no = ?", mobileNo).
		Where("password = ?", hex.EncodeToString(hash[:])).
		First(&Member).Error
	if err != nil {
		return nil, err
	}
	return Member, nil
}

type RiderMainRegisterStruct struct {
	Id       int    `gorm:"primary_key" json:"id"`
	ShardKey int    `gorm:"column:shard_key" json:"shard_key"`
	MobileNo string `gorm:"column:mobile_no" json:"mobile_no"`
}

type RiderRegisterStruct struct {
	Id                int    `gorm:"primary_key" json:"id"`
	Username          string `gorm:"column:username" json:"username"`
	MobileNo          string `gorm:"column:mobile_no" json:"mobile_no"`
	Email             string `gorm:"column:email" json:"email"`
	Password          string `gorm:"column:password" json:"password"`
	SecondaryPassword string `gorm:"column:secondary_password" json:"secondary_password"`
}

func RiderRegister(tx *gorm.DB, username string, email string, pass string, mobileNo string) (*RiderRegisterStruct, error) {
	//idx := helpers.ShardHash(mobileNo)
	hash := sha256.Sum256([]byte(pass))
	//key, _ := strconv.Atoi(idx)
	var checkUser RiderLoginStruct
	err := tx.Table("riders").
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

	registrationClone := RiderRegisterStruct{
		Username:          username,
		Email:             email,
		MobileNo:          mobileNo,
		Password:          hex.EncodeToString(hash[:]),
		SecondaryPassword: hex.EncodeToString(hash[:]),
	}
	err = tx.Table("riders").Create(&registrationClone).Error
	if err != nil {
		return nil, err
	}
	return &registrationClone, nil
}
