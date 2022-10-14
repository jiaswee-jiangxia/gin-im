package model

type UserStruct struct {
	BaseModel
	Username string `gorm:"column:username" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	Contact  string `gorm:"column:contact" json:"contact"`
}

func GetProfile(username string) (*UserStruct, error) {
	var Member *UserStruct
	//idx := helpers.ShardHash(username)
	err := db.Table("users").
		Where("username = ?", username).
		First(&Member).Error
	if err != nil {
		return nil, err
	}
	return Member, nil
}
