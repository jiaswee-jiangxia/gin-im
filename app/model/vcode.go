package model

type Vcode struct {
	BaseModel
	Cred       string `form:"cred" json:"cred"`
	CredType   string `form:"cred_type" json:"cred_type"`
	Vcode      string `form:"vcode" json:"vcode"`
	ExpiryTime int64  `form:"expiry" json:"expiry"`
	Purpose    string `form:"purpose" json:"purpose"`
}

func SaveVcode(cred, cred_type, vcode, purpose string, expiry int64) error {
	vcodeObj := Vcode{
		Cred:       cred,
		CredType:   cred_type,
		Vcode:      vcode,
		Purpose:    purpose,
		ExpiryTime: expiry,
	}
	err := db.Table("vcode").Create(&vcodeObj).Error
	return err
}
