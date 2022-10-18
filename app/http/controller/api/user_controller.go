package api

import (
	"encoding/json"
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/service/redis_service"
	"goskeleton/app/service/user_service"
	"goskeleton/app/utils/response"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(context *gin.Context) {
	var creds Credentials
	if err := context.ShouldBindJSON(&creds); err != nil {
		response.ErrorParam(context, creds)
		return
	}
	expirationTime := time.Now().Add(720 * time.Hour)
	userService := user_service.TokenStruct{
		Username: creds.Username,
		Password: creds.Password,
	}
	member, err := userService.UserLogin()
	if err != nil || member.Id <= 0 {
		response.SuccessButFail(context, "invalid_username_password", "invalid_username_password", nil)
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
		response.SuccessButFail(context, errSignedString.Error(), "ok", nil)
		return
	}
	response.Success(context, "ok", signedString)
	return
}

type RegisterStruct struct {
	Username             string `json:"username" binding:"required"`
	Password             string `json:"password" binding:"required"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required"`
	Email                string `json:"email"`
	Contact              string `json:"contact" binding:"required"`
}

func Register(context *gin.Context) {
	var creds RegisterStruct
	if err := context.ShouldBindJSON(&creds); err != nil {
		response.ErrorParam(context, creds)
		return
	}
	if creds.Password != creds.ConfirmationPassword {
		response.SuccessButFail(context, "wrong_confirmation_password", "ok", nil)
		return
	}

	db := model.GetDB()
	// begin a transaction
	txUser := db.Begin()
	tx := db.Begin()
	expirationTime := time.Now().Add(720 * time.Minute)
	userService := user_service.TokenStruct{
		Username: creds.Username,
		Contact:  creds.Contact,
		Email:    creds.Email,
		Password: creds.Password,
		Tx:       txUser,
	}
	member, err := userService.UserRegister()
	if err != nil {
		txUser.Rollback()
		response.SuccessButFail(context, err.Error(), "ok", nil)
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
		txUser.Rollback()
		tx.Rollback()
		response.SuccessButFail(context, errSignedString.Error(), "ok", nil)
		return
	}
	txUser.Commit()
	tx.Commit()
	response.Success(context, "ok", signedString)
	return
}

func GetProfile(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, "ok", nil)
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
		response.Success(context, "ok", returnProfile)
		return
	}
	userService := user_service.TokenStruct{
		Username: usernameText,
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), "ok", nil)
		return
	}
	rdb.CacheValue = profile
	rdb.PrepareCacheWrite()
	response.Success(context, "ok", profile)
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
		response.Success(context, "ok", nil)
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
	if prof.Nickname != nil {
		err = userService.UpdateNickname(prof.Nickname)
		if err != nil {
			response.SuccessButFail(context, err.Error(), "ok", nil)
			return
		}
	}
	if prof.Nickname != nil {
		err = userService.UpdateEmail(prof.Email)
		if err != nil {
			response.SuccessButFail(context, err.Error(), "ok", nil)
			return
		}
	}
	if prof.Nickname != nil {
		err = userService.UpdateContact(prof.Contact)
		if err != nil {
			response.SuccessButFail(context, err.Error(), "ok", nil)
			return
		}
	}
	if prof.Nickname != nil {
		err = userService.UpdateBFVerified(prof.BFVerified)
		if err != nil {
			response.SuccessButFail(context, err.Error(), "ok", nil)
			return
		}
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), "ok", nil)
		return
	}
	rdb.CacheValue = profile
	rdb.PrepareCacheWrite()
	response.Success(context, "ok", profile)
	return
}

type UpdateTokenStruct struct {
	Wx_token  *string `json:"wx_token"`
	Ios_token *string `json:"ios_token"`
}

func UpdateToken(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, "ok", nil)
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
			response.SuccessButFail(context, err.Error(), "ok", nil)
			return
		}
	}
	if tokens.Ios_token != nil {
		err := userService.UpdateIosToken(tokens.Ios_token)
		if err != nil {
			response.SuccessButFail(context, err.Error(), "ok", nil)
			return
		}
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), "ok", nil)
		return
	}
	rdb.CacheValue = profile
	rdb.PrepareCacheWrite()
	response.Success(context, "ok", profile)
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
		response.Success(context, "ok", nil)
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
		response.SuccessButFail(context, "password different", "ok", nil)
		return
	}
	if err != nil {
		response.SuccessButFail(context, err.Error(), "ok", nil)
		return
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), "ok", nil)
		return
	}

	response.Success(context, "ok", profile)
	return
}
