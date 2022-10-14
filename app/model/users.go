package model

type UserStruct struct {
	Id        string `gorm:"column:id" json:"-"`
	Username  string `gorm:"column:username" json:"username"`
	Email     string `gorm:"column:email" json:"email"`
	Contact   string `gorm:"column:contact" json:"contact"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
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
