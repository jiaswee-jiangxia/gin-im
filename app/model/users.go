package model

import (
	"encoding/json"
	"errors"
	"fmt"
	consts "goskeleton/app/global/response"
	"goskeleton/app/helpers"
	"goskeleton/app/service/redis_service"
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

type OTP struct {
	BaseModel
	Cred       string `gorm:"column:cred" json:"cred"`
	OTP        string `gorm:"column:otp" json:"otp"`
	ExpiryTime int64  `gorm:"column:expiry_time" json:"expiry_time"`
	Purpose    string `gorm:"column:purpose" json:"purpose"`
}

func UserRegister(username string, email string, pass string, mobileNo string) (*RegisterStruct, error) {
	var checkUser *Users
	var err error
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	cacheData := rdb.PrepareCacheRead()
	if cacheData != "" {
		return nil, errors.New("username_is_used")
	} else {
		err = db.Table("users").
			Where("username", username).First(&checkUser).Error
	}
	if checkUser.Id > 0 {
		rdb.CacheValue = checkUser
		rdb.PrepareCacheWrite()
		return nil, errors.New("username_is_used")
	}
	hash := helpers.GetMD5Hash(pass)
	registrationClone := &RegisterStruct{
		Username: username,
		Email:    email,
		Contact:  mobileNo,
		Password: hash,
	}
	err = db.Table("users").Create(registrationClone).Error
	if err != nil {
		return nil, err
	}
	return registrationClone, nil
}

func UserLogin(username string, pass string) (*Users, error) {
	var member *Users
	var err error
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	cacheData := rdb.PrepareCacheRead()
	if cacheData != "" {
		err = json.Unmarshal([]byte(cacheData), &member)
		if err != nil {
			return nil, err
		}
		return member, nil
	}
	hash := helpers.GetMD5Hash(pass)
	err = db.Table("users").
		Where("username = ?", username).
		Where("password = ?", hash).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	rdb.CacheValue = member
	rdb.PrepareCacheWrite()
	return member, nil
}

func GetUserByUsername(username string) (*Users, error) {
	var u *Users
	var err error
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	cacheData := rdb.PrepareCacheRead()
	if cacheData != "" {
		err = json.Unmarshal([]byte(cacheData), &u)
	} else {
		err = db.Table("users").
			Where("username = ?", username).
			First(&u).Error
		if err == nil {
			rdb.CacheValue = u
			rdb.PrepareCacheWrite()
		}
	}
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
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + *username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	rdb.DelCache()
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
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + *username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	rdb.DelCache()
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
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + *username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	rdb.DelCache()
	return err
}

func UpdateWxToken(username *string, wxtoken *string) error {
	err := db.Table("users").
		Where("username = ?", username).
		Update("wx_token", wxtoken).
		Error
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + *username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	rdb.DelCache()
	return err
}

func UpdateIosToken(username *string, iostoken *string) error {
	err := db.Table("users").
		Where("username = ?", username).
		Update("ios_token", iostoken).
		Error
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + *username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	rdb.DelCache()
	return err
}

func UpdatePassword(username *string, oldPassword string, newPassword string) error {
	var Member *Users
	// Redis Lock
	lockFlag := redis_service.PrepareLockTrial(redis_service.RedisCacheLock, "UPDATE_PASSWORD:"+*username, nil, 60)
	if !lockFlag {
		return errors.New(consts.WaitingPreviousActionToBeCompleted)
	}
	hash := helpers.GetMD5Hash(oldPassword)
	newHash := helpers.GetMD5Hash(newPassword)
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
	if err != nil {
		return err
	}
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + *username,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	rdb.DelCache()
	// Redis Unlock
	redis_service.PrepareUnlockTrial(redis_service.RedisCacheLock, "UPDATE_PASSWORD")
	return err
}

func SaveOTP(cred, otp, purpose string, expiry int64) error {
	otpObj := OTP{
		Cred:       cred,
		OTP:        otp,
		Purpose:    purpose,
		ExpiryTime: expiry,
	}
	err := db.Table("otp").Create(&otpObj).Error
	return err
}

func UserLoginWithEmail(email string, otp string) (*Users, error) {
	var member = &Users{}
	var otpp = &OTP{}
	var err error
	db.Table("otp").
		Where("cred = ?", email).
		Where("otp = ?", otp).First(&otpp)
	if otpp.Id > 0 { // Record matched
		err = db.Table("users").
			Where("email", email).First(&member).Error
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(member)
	return member, nil
}
