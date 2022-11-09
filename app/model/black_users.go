package model

// BlackUsers struct
type BlackUsers struct {
	BaseModel
	UserId      int64 `gorm:"column:user_id" json:"user_id"`
	BlackUserId int64 `gorm:"column:black_user_id" json:"black_user_id"`
	Active      int64 `gorm:"column:active" json:"active"`
}

type BlackUsersUser struct {
	BlackUsers
	Username string `gorm:"column:username" json:"username"`
}

func CreateNewBlackUser(u *BlackUsers) (*BlackUsers, error) {
	err := db.Table("black_users").
		Create(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetBlackUser(userId string, blackUserId string) (*BlackUsers, error) {
	var u *BlackUsers
	err := db.Table("black_users").
		Where("user_id = ?", userId).
		Where("black_user_id = ?", blackUserId).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (g BlackUsers) Updates(updates interface{}) (*BlackUsers, error) {
	err := db.Table("black_users").Model(&g).
		Where("user_id = ?", g.UserId).
		Where("black_user_id = ?", g.BlackUserId).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func GetBlackUserList(u *BlackUsers, limit int, offset int) ([]*BlackUsersUser, error) {
	var g []*BlackUsersUser
	var err error
	err = db.Table("black_users").
		Joins("inner join users on users.id = black_users.black_user_id").
		Where("user_id = ?", u.UserId).
		Where("active = ?", 1).
		Select("black_users.*, users.username").
		Limit(limit).
		Offset(offset).
		Find(&g).Error
	if err != nil {
		return nil, err
	}
	return g, nil
}
