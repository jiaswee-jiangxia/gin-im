package api

import (
	consts "goskeleton/app/global/response"
	"goskeleton/app/helpers"
	"goskeleton/app/service/vcode_service"
	"goskeleton/app/utils/response"
	"net/smtp"

	"github.com/gin-gonic/gin"
)

const EXPIRY_TIME = 0

type VcodeRequest struct {
	Purpose string `form:"purpose" json:"purpose" binding:"required"` // Login, register, verify etc
	IdType  string `form:"id_type" json:"id_type" binding:"required"` // Phone, email, token etc
	Id      string `form:"id" json:"id" binding:"required"`           // +650000000, abc@def.com, etc
}

func GetVcode(context *gin.Context) {
	vcode := &vcode_service.Vcode{}
	Req := VcodeRequest{}
	if err := context.ShouldBind(&Req); err != nil {
		response.ErrorParam(context, Req)
		return
	}
	vcode.Purpose = Req.Purpose
	if Req.IdType == "email" { // Request for email Vcode
		EmailVcode(*vcode, Req)
	} else if Req.IdType == "phone" { // Request for phone Vcode
		PhoneVcode(*vcode, Req)
	}
	response.Success(context, consts.Success, "")
	return
}

func EmailVcode(vcode vcode_service.Vcode, Req VcodeRequest) {
	vcode.Cred = Req.Id
	vcode.Vcode = "000000" // Generate with OTP generator, hardcode for now
	vcode.CredType = "email"
	vcode.ExpiryTime = EXPIRY_TIME
	vcode.SaveVcode()

	from := "wuikian@jiangxia.com.sg" // Replace with sender email
	password := "Wuikian789!@#"       // Replace with sender email password
	toEmailAddress := Req.Id
	to := []string{toEmailAddress}

	host := "mail.jiangxia.com.sg" // Email host
	port := "587"                  // Email host port
	address := host + ":" + port

	subject := "Subject: This is the subject of the mail\n" // Email subject
	body := vcode.Vcode                                     // OTP code and other message
	message := []byte(subject + "\n" + body)

	auth := helpers.LoginAuthWrapper{
		Username: from,
		Password: password,
	}
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		panic(err)
	}
}
func PhoneVcode(vcode vcode_service.Vcode, Req VcodeRequest) {
	vcode.Cred = Req.Id
	vcode.CredType = "phone"
	vcode.Vcode = "000000" // Generate with OTP generator, hardcode for now
	vcode.ExpiryTime = EXPIRY_TIME
	vcode.SaveVcode()
	SendMessage(vcode)
}

func SendMessage(vcode vcode_service.Vcode) {
	// To be implement
	return
}
