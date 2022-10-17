package group_service

import (
	"goskeleton/app/model"
)

type GroupStruct struct {
	model.BaseModel
	Name      string `gorm:"column:name" json:"groupname"`
	CreatedBy string `gorm:"column:created_by" json:"created_by"`
	Owner     string `gorm:"column:owner" json:"owner"`
	Disbanded bool   `gorm:"column:disbanded" json:"disbanded"`
}

type GroupMemberStruct struct {
	model.BaseModel
	GroupID  int64  `gorm:"column:group_id" json:"group_id"`
	Username string `gorm:"column:username" json:"user_id"`
	Role     string `gorm:"column:role" json:"role"`
}

func (m *GroupStruct) CreateGroup() (*model.GroupStruct, error) {
	group, err := model.CreateGroup(m.Name, m.Owner)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (m *GroupStruct) GetGroupInfo() (*model.GroupStruct, error) {
	groupInfo, err := model.GetGroupInfo(m.Id)
	if err != nil {
		return nil, err
	}
	return groupInfo, nil
}

func (m *GroupStruct) GetGroupAdmin() ([]model.GroupMemberStruct, error) {
	adminInfo, err := model.GetGroupAdminInfo(m.Id)
	if err != nil {
		return nil, err
	}
	return adminInfo, nil
}

func (m *GroupStruct) GetGroupMember() ([]model.GroupMemberStruct, error) {
	adminInfo, err := model.GetGroupMemberInfo(m.Id)
	if err != nil {
		return nil, err
	}
	return adminInfo, nil
}

func (m *GroupStruct) AddGroupMember(username string) error {
	err := model.AddGroupMember(m.Id, username)
	if err != nil {
	}
	return nil
}

func (m *GroupStruct) SetGroupAdmin(memberUsername string) error {
	err := model.SetGroupAdmin(m.Id, memberUsername)
	if err != nil {
		return err
	}
	return nil
}

func (m *GroupStruct) SetGroupOwner(memberUsername string) error {
	err := model.SetGroupOwner(m.Id, memberUsername)
	if err != nil {
		return err
	}
	return nil
}

func (m *GroupStruct) RemoveGroupMember(memberUsename string) error {
	err := model.RemoveGroupMember(m.Id, memberUsename)
	if err != nil {
		return err
	}
	return nil
}

func (m *GroupStruct) DisbandGroup() error {
	err := model.DisbandGroup(m.Id)
	if err != nil {
		return err
	}
	return nil
}
