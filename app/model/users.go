package model

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
