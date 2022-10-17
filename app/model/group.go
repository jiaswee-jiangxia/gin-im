package model

type GroupStruct struct {
	BaseModel
	Name      string `gorm:"column:name" json:"groupname"`
	CreatedBy string `gorm:"column:created_by" json:"created_by"`
	Owner     string `gorm:"column:owner" json:"owner"`
	Disbanded bool   `gorm:"column:disbanded" json:"disbanded"`
}

type GroupMemberStruct struct {
	BaseModel
	GroupID  int64  `gorm:"column:group_id" json:"group_id"`
	Username string `gorm:"column:username" json:"user_id"`
	Role     string `gorm:"column:role" json:"role"`
}

func CreateGroup(groupName string, owner string) (*GroupStruct, error) {
	group := &GroupStruct{
		Name:      groupName,
		CreatedBy: owner,
		Owner:     owner,
		Disbanded: false,
	}
	err := db.Table("groups").Create(&group).Error
	if err != nil {
		return nil, err
	}
	return group, nil
}

func AddGroupMember(groupID int64, memberUsername string) error {
	err := db.Table("group_members").Select("group_id", "username", "role").
		Create(&GroupMemberStruct{
			GroupID:  groupID,
			Username: memberUsername,
			Role:     "member",
		}).Error
	return err
}

func GetGroupInfo(groupID int64) (*GroupStruct, error) {
	res := &GroupStruct{}
	err := db.Table("groups").Where("id", groupID).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetGroupAdminInfo(groupID int64) ([]GroupMemberStruct, error) {
	res := make([]GroupMemberStruct, 0)
	db.Table("group_members").Where("group_id", groupID).Where("role", "admin").Scan(&res)
	return res, nil
}

func GetGroupMemberInfo(groupID int64) ([]GroupMemberStruct, error) {
	res := make([]GroupMemberStruct, 0)
	db.Table("group_members").Where("group_id", groupID).Scan(&res)
	return res, nil
}

func SetGroupAdmin(groupID int64, memberUsername string) error {
	err := db.Table("group_members").Where("group_id", groupID).Where("username", memberUsername).Update("role", "admin").Error
	return err
}

func SetGroupOwner(groupID int64, memberUsername string) error {
	err := db.Table("groups").Where("id", groupID).Update("owner", memberUsername).Error
	return err
}

func RemoveGroupMember(groupID int64, memberUsername string) error {
	member := &GroupMemberStruct{}
	err := db.Table("group_members").Where("group_id", groupID).Where("username", memberUsername).Delete(&member).Error
	return err
}

func DisbandGroup(groupID int64) error {
	member := &GroupMemberStruct{}
	err := db.Table("group_members").Where("group_id", groupID).Delete(&member).Error
	if err != nil {
		return err
	}
	err = db.Table("groups").Where("id", groupID).Update("disbanded", true).Error
	return err
}
