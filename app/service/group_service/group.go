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

func (m *GroupStruct) GetGroupInfo() (g *model.GroupStruct, err error) {
	g, err = model.GetGroupInfo(m.Id)
	if err != nil {
		return nil, err
	}
	return
}

func (m *GroupStruct) GetGroupAdmin() (g []model.GroupMemberStruct, err error) {
	g, err = model.GetGroupAdminInfo(m.Id)
	if err != nil {
		return nil, err
	}
	return
}

func (m *GroupStruct) GetGroupMember() (g []model.GroupMemberStruct, err error) {
	g, err = model.GetGroupMemberInfo(m.Id)
	return
}

func (m *GroupStruct) AddGroupMember(username string) error {
	err := model.AddGroupMember(m.Id, username)
	return err
}

func (m *GroupStruct) SetGroupAdmin(memberUsername string) error {
	err := model.SetGroupAdmin(m.Id, memberUsername)
	return err
}

func (m *GroupStruct) SetGroupOwner(memberUsername string) error {
	err := model.SetGroupOwner(m.Id, memberUsername)
	return err
}

func (m *GroupStruct) RemoveGroupMember(memberUsename string) error {
	err := model.RemoveGroupMember(m.Id, memberUsename)
	return err
}

func (m *GroupStruct) DisbandGroup() error {
	err := model.DisbandGroup(m.Id)
	return err
}
