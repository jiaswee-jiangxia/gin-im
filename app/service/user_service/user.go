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

func (m *TokenStruct) UserLogin() (*model.Users, error) {
	member, err := model.UserLogin(m.Username, m.Password)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *TokenStruct) UserRegister() (*model.RegisterStruct, error) {
	member, err := model.UserRegister(m.Username, m.Email, m.Password, m.Contact)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *TokenStruct) UserProfile() (*model.Users, error) {
	profile, err := model.GetUserByUsername(m.Username)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (m *TokenStruct) FindUserByUsername() (*model.Users, error) {
	u, err := model.GetUserByUsername(m.Username)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *TokenStruct) UpdateNickname(nickname *string) error {
	err := model.UpdateNickname(&m.Username, nickname)
	return err
}

func (m *TokenStruct) UpdateContact(contact *string) error {
	err := model.UpdateContact(&m.Username, contact)
	return err
}

func (m *TokenStruct) UpdateEmail(email *string) error {
	err := model.UpdateEmail(&m.Username, email)
	return err
}

func (m *TokenStruct) UpdateBFVerified(BFVerified *bool) error {
	err := model.UpdateBFVerified(&m.Username, BFVerified)
	return err
}

func (m *TokenStruct) UpdateWxToken(WxToken *string) error {
	err := model.UpdateWxToken(&m.Username, WxToken)
	return err
}

func (m *TokenStruct) UpdateIosToken(IosToken *string) error {
	err := model.UpdateIosToken(&m.Username, IosToken)
	return err
}

func (m *TokenStruct) UpdatePassword(old_password string, new_password string) error {
	err := model.UpdatePassword(&m.Username, old_password, new_password)
	return err
}
