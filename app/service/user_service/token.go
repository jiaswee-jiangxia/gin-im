package user_service

import (
	"gorm.io/gorm"
	"goskeleton/app/model"
)

type TokenStruct struct {
	Username string
	Email    string
	Contact  string
	Password string
	UserId   string
	Tx       *gorm.DB
}

func (m *TokenStruct) UserLogin() (*model.LoginStruct, error) {
	member, err := model.UserLogin(m.Username, m.Password)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *TokenStruct) UserRegister() (*model.RegisterStruct, error) {
	member, err := model.UserRegister(m.Tx, m.Username, m.Email, m.Password, m.Contact)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *TokenStruct) UserProfile() (*model.UserStruct, error) {
	profile, err := model.GetProfile(m.Username)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
