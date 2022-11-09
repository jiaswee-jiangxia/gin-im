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
	Type     string `form:"type" json:"type" binding:"required"`
	Passcode string `form:"passcode" json:"passcode"`
	Id       string `form:"id" json:"id" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}

func Login(context *gin.Context) {
	var creds Credentials
	if err := context.ShouldBind(&creds); err != nil {
		response.ErrorParam(context, creds)
		return
	}
	if creds.Type == "password" {
		LoginByPassword(context, creds)
		return
	}
	if creds.Type == "email" {
		LoginByEmail(context, creds)
		return
	}
}

func LoginByPassword(context *gin.Context, creds Credentials) {
	expirationTime := time.Now().Add(720 * time.Hour)
	userService := user_service.TokenStruct{
		Username: creds.Id,
		Password: &creds.Passcode,
	}
	hash := helpers.GetMD5Hash(creds.Passcode)
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
	response.Success(context, consts.Success, &Token{Token: signedString})
	return
}

func LoginByEmail(context *gin.Context, creds Credentials) {
	expirationTime := time.Now().Add(720 * time.Hour)
	userService := user_service.TokenStruct{
		Email: &creds.Id,
	}
	member, err := userService.UserLoginWithEmail(creds.Id)
	if err != nil || member.Id <= 0 {

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
	response.Success(context, consts.Success, &Token{Token: signedString})
	return
}

type RegisterStruct struct {
	Method               string    `form:"method" json:"method" binding:"required"`
	Username             string    `form:"username" json:"username" binding:"required,alphanum,min=4"`
	Password             string    `form:"password" json:"password"`
	ConfirmationPassword string    `form:"confirmation_password" json:"confirmation_password"`
	Email                *string   `form:"email" json:"email,omitempty" binding:"omitempty,email"`
	Contact              string    `form:"contact" json:"contact"`
	Vcode                string    `form:"vcode" json:"vcode"`
	PhoneCode            PhoneCode `form:"phone_code" json:"phone_code"`
}

type PhoneCode struct {
	Country     string `form:"country" json:"country"`
	Code        string `form:"code" json:"code"`
	CountryFull string `form:"country_full" json:"country_full"`
}

func Register(context *gin.Context) {
	var creds RegisterStruct
	if err := context.ShouldBindJSON(&creds); err != nil {
		response.ErrorParam(context, creds)
		return
	}
	expirationTime := time.Now().Add(720 * time.Minute)
	userService := user_service.TokenStruct{}

	switch method := creds.Method; method {
	case "password":
		if creds.Password != creds.ConfirmationPassword {
			response.SuccessButFail(context, consts.WrongConfirmationPassword, consts.Failed, nil)
			return
		}
		userService = user_service.TokenStruct{
			Username: creds.Username,
			Password: &creds.Password,
		}
	case "email":
		// TODO:match vcode
		// TODO:check duplicate email
		userService = user_service.TokenStruct{
			Email: creds.Email,
			Vcode: creds.Vcode,
		}
	case "phone":
		// TODO:match vcode
		// TODO:check duplicate phone
		userService = user_service.TokenStruct{
			Contact:      &creds.Contact,
			PhoneCountry: &creds.PhoneCode.Country,
			PhoneCode:    creds.PhoneCode.Code,
			CountryFull:  creds.PhoneCode.CountryFull,
			Vcode:        creds.Vcode,
		}
	default:
		fmt.Println("Error")
		return
	}
	fmt.Println(userService.Email)
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
	response.Success(context, consts.Success, &Token{Token: signedString})
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
	lockFlag := redis_service.PrepareLockTrial(redis_service.RedisCacheLock, "UPDATE_PROFILE:"+usernameText, nil, 60)
	if !lockFlag {
		response.SuccessButFail(context, consts.WaitingPreviousActionToBeCompleted, consts.WaitingPreviousActionToBeCompleted, nil)
		return
	}
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
	OldPassword     string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword     string `form:"new_password" json:"new_password" binding:"required"`
	ConfirmPassword string `form:"confirmation_password" json:"confirmation_password" binding:"required"`
}

func UpdatePassword(context *gin.Context) {
	var err error
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)

	var passwords UpdatePasswordStruct
	if err := context.ShouldBind(&passwords); err != nil { // Get request data
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

func GetPhoneNo(context *gin.Context) {
	var err error
  username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)
  
	userService := user_service.TokenStruct{
		Username: usernameText,
	}
	profile, err := userService.UserProfile()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
  
	bytes := []byte(`{"PhoneNo":` + profile.Contact + "}")
	var data map[string]interface{} = make(map[string]interface{})
	json.Unmarshal(bytes, &data)
	response.Success(context, consts.Success, data)
}


func RefreshToken(context *gin.Context) {
  username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	expirationTime := time.Now().Add(720 * time.Minute)

	// Get Profile
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_PROFILE:" + usernameText,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	cacheData := rdb.PrepareCacheRead()
	if cacheData != "" {
		var returnProfile interface{}
		_ = json.Unmarshal([]byte(cacheData), &returnProfile)
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
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: usernameText,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        strconv.Itoa(int(profile.Id)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedString, errSignedString := token.SignedString(variable.PrivateKey)
	if errSignedString != nil {
		response.SuccessButFail(context, errSignedString.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, &Token{Token: signedString})
	return
}

func CheckToken(context *gin.Context) {
	response.Success(context, consts.Success, nil)
	return
}
