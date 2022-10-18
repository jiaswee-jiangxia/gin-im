package model

import (
	"errors"
	"goskeleton/app/helpers"

	"gorm.io/gorm"
)

type Users struct {
	BaseModel
	Username   string `gorm:"column:username" json:"username"`
	Password   string `gorm:"column:password" json:"password"`
	Nickname   string `gorm:"column:nickname" json:"nickname"`
	Email      string `gorm:"column:email" json:"email"`
	Contact    string `gorm:"column:contact" json:"contact"`
	BFVerified int64  `gorm:"column:b_f_verified" json:"b_f_verified"`
	WxToken    string `gorm:"column:wx_token" json:"wx_token"`
	IosToken   string `gorm:"column:ios_token" json:"ios_token"`
	LastOnline string `gorm:"column:last_online" json:"last_online"`
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
	if checkUser.Id > 0 {
		return nil, errors.New("username_is_used")
	}
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

func GetUserByUsername(username string) (*Users, error) {
	var u *Users
	err := db.Table("users").
		Where("username = ?", username).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UpdateNickname(username *string, nickname *string) error {
	var err error
	if nickname != nil {
		err = db.Table("users").
			Where("username = ?", username).
			Update("nickname", nickname).
			Error
	}
	return err
}

func UpdateEmail(username *string, email *string) error {
	var err error
	if email != nil {
		err = db.Table("users").
			Where("username = ?", username).
			Update("email", email).
			Error
	}
	return err
}

func UpdateContact(username *string, contact *string) error {
	var err error
	if contact != nil {
		err = db.Table("users").
			Where("username = ?", username).
			Update("contact", contact).
			Error
	}
	return err
}

func UpdateBFVerified(username *string, BFVerified *bool) error {
	var err error
	if BFVerified != nil {
		err = db.Table("users").
			Where("username = ?", username).
			Update("b_f_verified", BFVerified).
			Error
	}
	return err
}

func UpdateWxToken(username *string, wxtoken *string) error {
	err := db.Table("users").
		Where("username = ?", username).
		Update("wx_token", wxtoken).
		Error
	return err
}

func UpdateIosToken(username *string, iostoken *string) error {
	err := db.Table("users").
		Where("username = ?", username).
		Update("ios_token", iostoken).
		Error
	return err
}

func UpdatePassword(username *string, old_password string, new_password string) error {
	var Member *Users
	hash := helpers.GetMD5Hash(old_password)
	newHash := helpers.GetMD5Hash(new_password)
	err := db.Table("users").
		Where("username = ?", username).
		First(&Member).Error
	if err != nil {
		return err
	}
	if Member.Password == hash {
		err = db.Table("users").
			Where("username = ?", username).
			Update("password", newHash).
			Error
	}

	return err
}
