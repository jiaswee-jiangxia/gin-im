package api

import (
	"encoding/json"
	"fmt"
	consts "goskeleton/app/global/response"
	"goskeleton/app/helpers"

	"goskeleton/app/global/variable"
	"goskeleton/app/service/redis_service"
	"goskeleton/app/service/user_service"
	"goskeleton/app/utils/response"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Password string `form:"password" json:"password" binding:"required,alphanum,min=4"`
	Username string `form:"username" json:"username" binding:"required,min=6"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(context *gin.Context) {
	var creds Credentials
	if err := context.ShouldBind(&creds); err != nil {
		response.ErrorParam(context, creds)
		return
	}
	expirationTime := time.Now().Add(720 * time.Hour)
	userService := user_service.TokenStruct{
		Username: creds.Username,
		Password: creds.Password,
	}
	hash := helpers.GetMD5Hash(creds.Password)
	member, err := userService.UserLogin()
	if err != nil || member.Password != hash {
		response.SuccessButFail(context, consts.InvalidUsernamePassword, consts.InvalidUsernamePassword, nil)
		return
	}
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: member.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        strconv.FormatInt(member.Id, 10),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedString, errSignedString := token.SignedString(variable.PrivateKey)
	if errSignedString != nil {
		response.SuccessButFail(context, errSignedString.Error(), consts.Success, nil)
		return
	}
	response.Success(context, consts.Success, signedString)
	return
}

type RegisterStruct struct {
	Username             string `form:"username" json:"username" binding:"required,alphanum,min=4"`
	Password             string `form:"password" json:"password" binding:"required,min=6"`
	ConfirmationPassword string `form:"confirmation_password" json:"confirmation_password" binding:"required,min=6"`
	Email                string `form:"email" json:"email" binding:"email"`
	Contact              string `form:"contact" json:"contact" binding:"required,min=10"`
}

func Register(context *gin.Context) {
	var creds RegisterStruct
	if err := context.ShouldBindJSON(&creds); err != nil {
		response.ErrorParam(context, creds)
		return
	}
	if creds.Password != creds.ConfirmationPassword {
		response.SuccessButFail(context, consts.WrongConfirmationPassword, consts.Failed, nil)
		return
	}

	expirationTime := time.Now().Add(720 * time.Minute)
	userService := user_service.TokenStruct{
		Username: creds.Username,
		Contact:  creds.Contact,
		Email:    creds.Email,
		Password: creds.Password,
	}
	member, err := userService.UserRegister()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: member.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        strconv.Itoa(int(member.Id)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedString, errSignedString := token.SignedString(variable.PrivateKey)
	if errSignedString != nil {
		response.SuccessButFail(context, errSignedString.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, signedString)
	return
}

func GetProfile(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + usernameText,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	cacheData := rdb.PrepareCacheRead()
	if cacheData != "" {
		var returnProfile interface{}
		_ = json.Unmarshal([]byte(cacheData), &returnProfile)
		response.Success(context, consts.Success, returnProfile)
		return
	}
	userService := user_service.TokenStruct{
		Username: usernameText,
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	rdb.CacheValue = profile
	rdb.PrepareCacheWrite()
	response.Success(context, consts.Success, profile)
	return
}

type UpdateProfileStruct struct {
	Nickname   *string `json:"nickname"`
	Email      *string `json:"email"`
	Contact    *string `json:"contact"`
	BFVerified *bool   `json:"b_f_verified"`
}

func UpdateProfile(context *gin.Context) {
	var err error
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)

	var prof UpdateProfileStruct
	if err := context.ShouldBindJSON(&prof); err != nil { // Get request data
		response.ErrorParam(context, prof)
		return
	}
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + usernameText,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	userService := user_service.TokenStruct{
		Username: usernameText,
	}
	// Redis Lock
	redis_service.PrepareLockTrial(redis_service.RedisCacheLock, "UPDATE_PROFILE:"+usernameText, nil, 60)
	time.Sleep(10 * time.Second)
	if prof.Nickname != nil {
		err = userService.UpdateNickname(prof.Nickname)
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.Failed, nil)
			return
		}
	}
	if prof.Nickname != nil {
		err = userService.UpdateEmail(prof.Email)
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.Failed, nil)
			return
		}
	}
	if prof.Nickname != nil {
		err = userService.UpdateContact(prof.Contact)
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.Failed, nil)
			return
		}
	}
	if prof.Nickname != nil {
		err = userService.UpdateBFVerified(prof.BFVerified)
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.Failed, nil)
			return
		}
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	rdb.CacheValue = profile
	rdb.PrepareCacheWrite()
	// Redis UnLock
	redis_service.PrepareUnlockTrial(redis_service.RedisCacheLock, "UPDATE_PROFILE:"+usernameText)
	response.Success(context, consts.Success, profile)
	return
}

type UpdateTokenStruct struct {
	Wx_token  *string `json:"wx_token"`
	Ios_token *string `json:"ios_token"`
}

func UpdateToken(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Success, nil)
	}
	usernameText := fmt.Sprintf("%v", username)

	var tokens UpdateTokenStruct
	if err := context.ShouldBindJSON(&tokens); err != nil { // Get request data
		response.ErrorParam(context, tokens)
		return
	}
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + usernameText,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	userService := user_service.TokenStruct{
		Username: usernameText,
	}
	if tokens.Wx_token != nil {
		err := userService.UpdateWxToken(tokens.Wx_token)
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.Failed, nil)
			return
		}
	}
	if tokens.Ios_token != nil {
		err := userService.UpdateIosToken(tokens.Ios_token)
		if err != nil {
			response.SuccessButFail(context, err.Error(), consts.Failed, nil)
			return
		}
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	rdb.CacheValue = profile
	rdb.PrepareCacheWrite()
	response.Success(context, consts.Success, profile)
	return
}

type UpdatePasswordStruct struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirmation_password" binding:"required"`
}

func UpdatePassword(context *gin.Context) {
	var err error
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)

	var passwords UpdatePasswordStruct
	if err := context.ShouldBindJSON(&passwords); err != nil { // Get request data
		response.ErrorParam(context, passwords)
		return
	}

	userService := user_service.TokenStruct{
		Username: usernameText,
	}
	if passwords.NewPassword == passwords.ConfirmPassword {
		err = userService.UpdatePassword(passwords.OldPassword, passwords.NewPassword)
	} else {
		response.SuccessButFail(context, consts.WrongConfirmationPassword, consts.Failed, nil)
		return
	}
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}

	response.Success(context, consts.Success, profile)
	return
}
