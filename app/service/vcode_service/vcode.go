package vcode_service

import (
	"goskeleton/app/model"
)

type Vcode struct {
	Cred       string `form:"cred" json:"cred"`
	CredType   string `form:"cred_type" json:"cred_type"`
	Vcode      string `form:"vcode" json:"vcode"`
	ExpiryTime int64  `form:"expiry" json:"expiry"`
	Purpose    string `form:"purpose" json:"purpose"`
}

func (o *Vcode) SaveVcode() error {
	err := model.SaveVcode(o.Cred, o.CredType, o.Vcode, o.Purpose, o.ExpiryTime)
	return err
}
