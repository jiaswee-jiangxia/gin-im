package model

// Contacts struct
type Contacts struct {
	BaseModel
	UserId   int64  `gorm:"column:user_id" json:"user_id"`
	FriendId int64  `gorm:"column:friend_id" json:"friend_id"`
	Status   int64  `gorm:"column:status" json:"status"`
	Grouping string `gorm:"column:grouping" json:"grouping"`
}

func GetContactsByBothId(userId string, friendId string) (*Contacts, error) {
	var u *Contacts
	err := db.Table("contacts").
		Where("user_id = ?", userId).
		Where("friend_id = ?", friendId).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateNewContact(u *Contacts) (*Contacts, error) {
	err := db.Table("contacts").
		Create(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func Updates(g *Contacts, updates interface{}) (*Contacts, error) {
	err := db.Table("contacts").Model(&g).
		Where("user_id = ?", g.UserId).
		Where("friend_id = ?", g.FriendId).
		Updates(updates).Debug().Error
	if err != nil {
		return nil, err
	}
	return g, nil
}
